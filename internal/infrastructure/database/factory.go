package database

import (
	"fmt"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/database/adapters"
)

// AdapterFactory is a function that creates a new adapter instance
type AdapterFactory func() adapters.Adapter

var adapterRegistry = make(map[adapters.DriverType]AdapterFactory)

// RegisterAdapter registers a new database adapter
// This should be called from adapter init() functions
func RegisterAdapter(driverType adapters.DriverType, factory AdapterFactory) {
	adapterRegistry[driverType] = factory
}

// NewAdapter creates a new database adapter based on the driver type
func NewAdapter(driverType adapters.DriverType) (adapters.Adapter, error) {
	factory, exists := adapterRegistry[driverType]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedDriver, driverType)
	}
	return factory(), nil
}
