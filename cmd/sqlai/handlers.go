package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alessandrolattao/sqlai/internal/features/update"
)

// handleVersion prints the version information
func handleVersion() {
	fmt.Printf("SQL AI v%s\n", Version)
}

// handleUpdate checks for updates and installs them if available
func handleUpdate() {
	ctx := context.Background()
	updateService := update.NewService(Version)

	fmt.Printf("Current version: %s\n", Version)
	fmt.Println("Checking for updates...")

	latestVersion, hasUpdate, err := updateService.CheckForUpdate(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	if !hasUpdate {
		fmt.Printf("You are already running the latest version (%s)\n", Version)
		return
	}

	fmt.Printf("New version available: %s\n", latestVersion)
	fmt.Print("Do you want to update? [y/N]: ")

	var response string
	_, scanErr := fmt.Scanln(&response)
	if scanErr != nil && scanErr.Error() != "unexpected newline" {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", scanErr)
		os.Exit(1)
	}

	if response != "y" && response != "Y" {
		fmt.Println("Update cancelled")
		return
	}

	fmt.Println("Downloading and installing update...")
	if err := updateService.Update(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully updated to version %s!\n", latestVersion)
	fmt.Println("Please restart the application to use the new version.")
}
