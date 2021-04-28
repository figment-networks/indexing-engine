package datalake

import "strconv"

// DataLake represents raw data storage
type DataLake struct {
	network string
	chain   string

	storage Storage
}

// NewDataLake creates a data lake with the given storage provider
func NewDataLake(network string, chain string, storage Storage) *DataLake {
	return &DataLake{
		network: network,
		chain:   chain,
		storage: storage,
	}
}

// StoreResource stores the resource data
func (dl *DataLake) StoreResource(res *Resource) error {
	path := dl.resourcePath(res.Name)

	return dl.storage.Store(res.Data, path...)
}

// IsResourceStored checks if the resource is stored
func (dl *DataLake) IsResourceStored(name string) (bool, error) {
	path := dl.resourcePath(name)

	return dl.storage.IsStored(path...)
}

// RetrieveResource retrieves the resource data
func (dl *DataLake) RetrieveResource(name string) (*Resource, error) {
	path := dl.resourcePath(name)

	data, err := dl.storage.Retrieve(path...)
	if err != nil {
		return nil, err
	}

	return &Resource{
		Name: name,
		Data: data,
	}, nil
}

func (dl *DataLake) resourcePath(name string) []string {
	return []string{dl.network, dl.chain, name}
}

// StoreResourceAtHeight stores the resource data at the given height
func (dl *DataLake) StoreResourceAtHeight(res *Resource, height int64) error {
	path := dl.resourceAtHeightPath(res.Name, height)

	return dl.storage.Store(res.Data, path...)
}

// IsResourceStoredAtHeight checks if the resource is stored at the given height
func (dl *DataLake) IsResourceStoredAtHeight(name string, height int64) (bool, error) {
	path := dl.resourceAtHeightPath(name, height)

	return dl.storage.IsStored(path...)
}

// RetrieveResourceAtHeight retrieves the resource data at the given height
func (dl *DataLake) RetrieveResourceAtHeight(name string, height int64) (*Resource, error) {
	path := dl.resourceAtHeightPath(name, height)

	data, err := dl.storage.Retrieve(path...)
	if err != nil {
		return nil, err
	}

	return &Resource{
		Name: name,
		Data: data,
	}, nil
}

func (dl *DataLake) resourceAtHeightPath(name string, height int64) []string {
	h := strconv.FormatInt(height, 10)
	return []string{dl.network, dl.chain, "height", h, name}
}
