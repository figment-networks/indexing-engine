package datalake

// Storage is an interface for storing and retrieving raw data
type Storage interface {
	Store(data []byte, path ...string) error
	IsStored(path ...string) (bool, error)
	Retrieve(path ...string) ([]byte, error)
}
