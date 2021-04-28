package datalake

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFileStorage_ExistingDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "storage-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	_, err = NewFileStorage(dir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dir)

	assert.Nil(t, err)
	assert.True(t, info.IsDir())
}

func TestNewFileStorage_NonExistentDir(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "storage-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dir := filepath.Join(tmpDir, "non-existent")

	_, err = NewFileStorage(dir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dir)

	assert.Nil(t, err)
	assert.True(t, info.IsDir())
}

func TestFileStorage_Store(t *testing.T) {
	dir, err := ioutil.TempDir("", "storage-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fe, err := NewFileStorage(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = fe.Store([]byte("example data"), "path", "to", "file")
	if err != nil {
		t.Fatal(err)
	}

	file, err := ioutil.ReadFile(filepath.Join(dir, "path", "to", "file"))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "example data", string(file))
}

func TestFileStorage_IsStored(t *testing.T) {
	dir, err := ioutil.TempDir("", "storage-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fe, err := NewFileStorage(dir)
	if err != nil {
		t.Fatal(err)
	}

	stored, err := fe.IsStored("path", "to", "file")

	assert.NotNil(t, err)
	assert.False(t, stored)

	err = fe.Store([]byte{}, "path", "to", "file")
	if err != nil {
		t.Fatal(err)
	}

	stored, err = fe.IsStored("path", "to", "file")
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, stored)
}

func TestFileStorage_Retrieve_ExistingFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "storage-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fe, err := NewFileStorage(dir)
	if err != nil {
		t.Fatal(err)
	}

	err = fe.Store([]byte("example data"), "path", "to", "file")
	if err != nil {
		t.Fatal(err)
	}

	data, err := fe.Retrieve("path", "to", "file")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "example data", string(data))
}

func TestFileStorage_Retrieve_NonExistentFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "storage-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fe, err := NewFileStorage(dir)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fe.Retrieve("path", "to", "file")

	assert.NotNil(t, err)
}
