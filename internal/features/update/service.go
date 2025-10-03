// Package update provides self-update functionality for the sqlai binary.
package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubRepo = "alessandrolattao/sqlai"
	apiURL     = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
	timeout    = 30 * time.Second
)

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// Service handles binary self-update functionality
type Service struct {
	currentVersion string
}

// NewService creates a new update service
func NewService(currentVersion string) *Service {
	return &Service{
		currentVersion: currentVersion,
	}
}

// CheckForUpdate checks if a newer version is available on GitHub
func (s *Service) CheckForUpdate(ctx context.Context) (version string, hasUpdate bool, err error) {
	client := &http.Client{Timeout: timeout}

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if reqErr != nil {
		return "", false, fmt.Errorf("failed to create request: %w", reqErr)
	}

	resp, doErr := client.Do(req)
	if doErr != nil {
		return "", false, fmt.Errorf("failed to fetch release info: %w", doErr)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close response body: %w", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if decodeErr := json.NewDecoder(resp.Body).Decode(&release); decodeErr != nil {
		return "", false, fmt.Errorf("failed to parse release info: %w", decodeErr)
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	hasUpdate = compareVersions(s.currentVersion, latestVersion) < 0

	return latestVersion, hasUpdate, nil
}

// Update downloads and installs the latest version
func (s *Service) Update(ctx context.Context) (err error) {
	// Get latest release info
	client := &http.Client{Timeout: timeout}

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if reqErr != nil {
		return fmt.Errorf("failed to create request: %w", reqErr)
	}

	resp, doErr := client.Do(req)
	if doErr != nil {
		return fmt.Errorf("failed to fetch release info: %w", doErr)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close response body: %w", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if decodeErr := json.NewDecoder(resp.Body).Decode(&release); decodeErr != nil {
		return fmt.Errorf("failed to parse release info: %w", decodeErr)
	}

	// Find the correct asset for current platform
	assetName := getAssetName(strings.TrimPrefix(release.TagName, "v"))
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no release found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Download the asset
	tmpFile, downloadErr := downloadFile(ctx, downloadURL)
	if downloadErr != nil {
		return fmt.Errorf("failed to download update: %w", downloadErr)
	}
	defer func() {
		if rerr := os.Remove(tmpFile); rerr != nil && err == nil {
			err = fmt.Errorf("failed to remove temp download file: %w", rerr)
		}
	}()

	// Extract binary
	binaryPath, extractErr := extractBinary(tmpFile, release.TagName)
	if extractErr != nil {
		return fmt.Errorf("failed to extract binary: %w", extractErr)
	}
	defer func() {
		if rerr := os.Remove(binaryPath); rerr != nil && err == nil {
			err = fmt.Errorf("failed to remove extracted binary: %w", rerr)
		}
	}()

	// Replace current executable
	if replaceErr := replaceBinary(binaryPath); replaceErr != nil {
		return fmt.Errorf("failed to replace binary: %w", replaceErr)
	}

	return nil
}

// getAssetName returns the asset name for the current platform
func getAssetName(version string) string {
	var ext string
	if runtime.GOOS == "windows" {
		ext = ".zip"
	} else {
		ext = ".tar.gz"
	}

	return fmt.Sprintf("sqlai-%s-%s-%s%s", runtime.GOOS, runtime.GOARCH, version, ext)
}

// downloadFile downloads a file from a URL and returns the path to the temporary file
func downloadFile(ctx context.Context, url string) (tmpPath string, err error) {
	client := &http.Client{Timeout: 5 * time.Minute}

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if reqErr != nil {
		return "", fmt.Errorf("failed to create download request: %w", reqErr)
	}

	resp, doErr := client.Do(req)
	if doErr != nil {
		return "", fmt.Errorf("failed to download file: %w", doErr)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close response body: %w", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	tmpFile, createErr := os.CreateTemp("", "sqlai-update-*")
	if createErr != nil {
		return "", fmt.Errorf("failed to create temp file: %w", createErr)
	}

	tmpFileName := tmpFile.Name()

	if _, copyErr := io.Copy(tmpFile, resp.Body); copyErr != nil {
		if closeErr := tmpFile.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close temp file during cleanup: %w (copy error: %v)", closeErr, copyErr)
		}
		if removeErr := os.Remove(tmpFileName); removeErr != nil && err == nil {
			err = fmt.Errorf("failed to remove temp file during cleanup: %w (copy error: %v)", removeErr, copyErr)
		}
		if err == nil {
			err = fmt.Errorf("failed to write downloaded data: %w", copyErr)
		}
		return "", err
	}

	if closeErr := tmpFile.Close(); closeErr != nil {
		if removeErr := os.Remove(tmpFileName); removeErr != nil {
			return "", fmt.Errorf("failed to remove temp file: %w (close error: %v)", removeErr, closeErr)
		}
		return "", fmt.Errorf("failed to close temp file: %w", closeErr)
	}

	return tmpFileName, nil
}

// extractBinary extracts the binary from the downloaded archive
func extractBinary(archivePath, version string) (string, error) {
	binaryName := fmt.Sprintf("sqlai-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
		return extractFromZip(archivePath, binaryName)
	}
	return extractFromTarGz(archivePath, binaryName)
}

// extractFromTarGz extracts a specific file from a tar.gz archive
func extractFromTarGz(archivePath, fileName string) (tmpPath string, err error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close archive file: %w", cerr)
		}
	}()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() {
		if cerr := gzr.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close gzip reader: %w", cerr)
		}
	}()

	tr := tar.NewReader(gzr)

	for {
		header, readErr := tr.Next()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return "", fmt.Errorf("failed to read tar header: %w", readErr)
		}

		if header.Name == fileName || filepath.Base(header.Name) == fileName {
			tmpFile, createErr := os.CreateTemp("", "sqlai-binary-*")
			if createErr != nil {
				return "", fmt.Errorf("failed to create temp file: %w", createErr)
			}

			tmpFileName := tmpFile.Name()

			if _, copyErr := io.Copy(tmpFile, tr); copyErr != nil {
				if closeErr := tmpFile.Close(); closeErr != nil {
					err = fmt.Errorf("failed to close temp file during cleanup: %w (original error: %v)", closeErr, copyErr)
				}
				if removeErr := os.Remove(tmpFileName); removeErr != nil && err == nil {
					err = fmt.Errorf("failed to remove temp file during cleanup: %w (original error: %v)", removeErr, copyErr)
				}
				if err == nil {
					err = fmt.Errorf("failed to copy binary data: %w", copyErr)
				}
				return "", err
			}

			if closeErr := tmpFile.Close(); closeErr != nil {
				if removeErr := os.Remove(tmpFileName); removeErr != nil {
					return "", fmt.Errorf("failed to remove temp file: %w (close error: %v)", removeErr, closeErr)
				}
				return "", fmt.Errorf("failed to close temp file: %w", closeErr)
			}

			// Make executable
			if chmodErr := os.Chmod(tmpFileName, 0755); chmodErr != nil {
				if removeErr := os.Remove(tmpFileName); removeErr != nil {
					return "", fmt.Errorf("failed to remove temp file: %w (chmod error: %v)", removeErr, chmodErr)
				}
				return "", fmt.Errorf("failed to set executable permissions: %w", chmodErr)
			}

			return tmpFileName, nil
		}
	}

	return "", fmt.Errorf("binary not found in archive")
}

// extractFromZip extracts a specific file from a zip archive
func extractFromZip(archivePath, fileName string) (tmpPath string, err error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer func() {
		if cerr := reader.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close zip reader: %w", cerr)
		}
	}()

	for _, file := range reader.File {
		if file.Name == fileName || filepath.Base(file.Name) == fileName {
			rc, openErr := file.Open()
			if openErr != nil {
				return "", fmt.Errorf("failed to open file from zip: %w", openErr)
			}

			tmpFile, createErr := os.CreateTemp("", "sqlai-binary-*.exe")
			if createErr != nil {
				if closeErr := rc.Close(); closeErr != nil {
					return "", fmt.Errorf("failed to close zip entry: %w (create error: %v)", closeErr, createErr)
				}
				return "", fmt.Errorf("failed to create temp file: %w", createErr)
			}

			tmpFileName := tmpFile.Name()

			if _, copyErr := io.Copy(tmpFile, rc); copyErr != nil {
				if closeErr := rc.Close(); closeErr != nil {
					err = fmt.Errorf("failed to close zip entry during cleanup: %w (copy error: %v)", closeErr, copyErr)
				}
				if closeErr := tmpFile.Close(); closeErr != nil && err == nil {
					err = fmt.Errorf("failed to close temp file during cleanup: %w (copy error: %v)", closeErr, copyErr)
				}
				if removeErr := os.Remove(tmpFileName); removeErr != nil && err == nil {
					err = fmt.Errorf("failed to remove temp file during cleanup: %w (copy error: %v)", removeErr, copyErr)
				}
				if err == nil {
					err = fmt.Errorf("failed to copy binary data: %w", copyErr)
				}
				return "", err
			}

			if closeErr := rc.Close(); closeErr != nil {
				if closeErr2 := tmpFile.Close(); closeErr2 != nil {
					err = fmt.Errorf("failed to close temp file: %w (zip entry close error: %v)", closeErr2, closeErr)
				}
				if removeErr := os.Remove(tmpFileName); removeErr != nil && err == nil {
					err = fmt.Errorf("failed to remove temp file: %w (zip entry close error: %v)", removeErr, closeErr)
				}
				if err == nil {
					err = fmt.Errorf("failed to close zip entry: %w", closeErr)
				}
				return "", err
			}

			if closeErr := tmpFile.Close(); closeErr != nil {
				if removeErr := os.Remove(tmpFileName); removeErr != nil {
					return "", fmt.Errorf("failed to remove temp file: %w (close error: %v)", removeErr, closeErr)
				}
				return "", fmt.Errorf("failed to close temp file: %w", closeErr)
			}

			return tmpFileName, nil
		}
	}

	return "", fmt.Errorf("binary not found in archive")
}

// replaceBinary replaces the current executable with the new one
func replaceBinary(newBinaryPath string) (err error) {
	// Get current executable path
	currentExe, exeErr := os.Executable()
	if exeErr != nil {
		return fmt.Errorf("failed to get current executable path: %w", exeErr)
	}

	// Resolve symlinks
	currentExe, symlinkErr := filepath.EvalSymlinks(currentExe)
	if symlinkErr != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", symlinkErr)
	}

	// On Windows, we need to rename the old binary first
	if runtime.GOOS == "windows" {
		oldExe := currentExe + ".old"
		if renameErr := os.Rename(currentExe, oldExe); renameErr != nil {
			return fmt.Errorf("failed to rename old binary: %w", renameErr)
		}
		defer func() {
			if rerr := os.Remove(oldExe); rerr != nil && err == nil {
				err = fmt.Errorf("failed to remove old binary: %w", rerr)
			}
		}()
	}

	// Copy new binary to current location
	src, openErr := os.Open(newBinaryPath)
	if openErr != nil {
		return fmt.Errorf("failed to open new binary: %w", openErr)
	}
	defer func() {
		if cerr := src.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close source file: %w", cerr)
		}
	}()

	dst, createErr := os.OpenFile(currentExe, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if createErr != nil {
		return fmt.Errorf("failed to create destination file: %w", createErr)
	}
	defer func() {
		if cerr := dst.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close destination file: %w", cerr)
		}
	}()

	if _, copyErr := io.Copy(dst, src); copyErr != nil {
		return fmt.Errorf("failed to copy binary: %w", copyErr)
	}

	return nil
}

// compareVersions compares two version strings (e.g., "1.2.3" vs "1.3.0")
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 int
		if i < len(parts1) {
			_, _ = fmt.Sscanf(parts1[i], "%d", &p1)
		}
		if i < len(parts2) {
			_, _ = fmt.Sscanf(parts2[i], "%d", &p2)
		}

		if p1 < p2 {
			return -1
		}
		if p1 > p2 {
			return 1
		}
	}

	return 0
}
