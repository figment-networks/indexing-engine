package datalake

import (
	"errors"
	"strconv"
)

// DataLake represents raw data storage
type DataLake struct {
	network string
	chain   string

	storage Storage
}

// ErrResourceNameRequired is returned when the resource name is an empty string
var ErrResourceNameRequired = errors.New("resource name is required")

// NewDataLake creates a data lake with the given storage provider
func NewDataLake(network string, chain string, storage Storage) *DataLake {
	return &DataLake{
		network: network,
		chain:   chain,
		storage: storage,
	}
}

// StoreResource stores the resource data
func (dl *DataLake) StoreResource(res *Resource, name string) error {
	path, err := dl.resourcePath(name)
	if err != nil {
		return err
	}

	return dl.storage.Store(res.Data, path...)
}

// IsResourceStored checks if the resource is stored
func (dl *DataLake) IsResourceStored(name string) (bool, error) {
	path, err := dl.resourcePath(name)
	if err != nil {
		return false, err
	}

	return dl.storage.IsStored(path...)
}

// RetrieveResource retrieves the resource data
func (dl *DataLake) RetrieveResource(name string) (*Resource, error) {
	path, err := dl.resourcePath(name)
	if err != nil {
		return nil, err
	}

	data, err := dl.storage.Retrieve(path...)
	if err != nil {
		return nil, err
	}

	return &Resource{Data: data}, nil
}

func (dl *DataLake) resourcePath(name string) ([]string, error) {
	if name == "" {
		return nil, ErrResourceNameRequired
	}

	path := []string{dl.network, dl.chain, name}

	return path, nil
}

// StoreResourceAtHeight stores the resource data at the given height
func (dl *DataLake) StoreResourceAtHeight(res *Resource, name string, height int64) error {
	path, err := dl.resourceAtHeightPath(name, height)
	if err != nil {
		return err
	}

	return dl.storage.Store(res.Data, path...)
}

// IsResourceStoredAtHeight checks if the resource is stored at the given height
func (dl *DataLake) IsResourceStoredAtHeight(name string, height int64) (bool, error) {
	path, err := dl.resourceAtHeightPath(name, height)
	if err != nil {
		return false, err
	}

	return dl.storage.IsStored(path...)
}

// RetrieveResourceAtHeight retrieves the resource data at the given height
func (dl *DataLake) RetrieveResourceAtHeight(name string, height int64) (*Resource, error) {
	path, err := dl.resourceAtHeightPath(name, height)
	if err != nil {
		return nil, err
	}

	data, err := dl.storage.Retrieve(path...)
	if err != nil {
		return nil, err
	}

	return &Resource{Data: data}, nil
}

func (dl *DataLake) resourceAtHeightPath(name string, height int64) ([]string, error) {
	if name == "" {
		return nil, ErrResourceNameRequired
	}

	h := strconv.FormatInt(height, 10)
	path := []string{dl.network, dl.chain, "height", h, name}

	return path, nil
}
