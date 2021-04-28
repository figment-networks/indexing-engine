package datalake

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type fileStorage struct {
	directory string
}

var _ Storage = (*fileStorage)(nil)

// NewFileStorage creates a filesystem storage
func NewFileStorage(dir string) (Storage, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &fileStorage{directory: dir}, nil
}

func (fe *fileStorage) Store(data []byte, path ...string) error {
	filename := fe.fileName(path)

	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (fe *fileStorage) IsStored(path ...string) (bool, error) {
	info, err := os.Stat(fe.fileName(path))
	if err != nil {
		return false, err
	}

	return !info.IsDir(), nil
}

func (fe *fileStorage) Retrieve(path ...string) ([]byte, error) {
	return ioutil.ReadFile(fe.fileName(path))
}

func (fe *fileStorage) fileName(path []string) string {
	parts := append([]string{fe.directory}, path...)

	return filepath.Join(parts...)
}
