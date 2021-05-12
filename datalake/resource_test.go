package datalake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResource(t *testing.T) {
	res := NewResource([]byte("example data"))

	assert.Equal(t, []byte("example data"), res.Data)
}

func TestNewJSONResource(t *testing.T) {
	data := map[string]interface{}{"example": "data"}

	res, err := NewJSONResource(data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `{"example":"data"}`, string(res.Data))
}

func TestNewBinaryResource(t *testing.T) {
	res, err := NewBinaryResource(42)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []byte{0x3, 0x4, 0x0, 0x54}, res.Data)
}

func TestNewBase64Resource(t *testing.T) {
	res, err := NewBase64Resource("example data")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "DwwADGV4YW1wbGUgZGF0YQ==", string(res.Data))
}

func TestResource_ScanJSON(t *testing.T) {
	res := Resource{
		Data: []byte(`{"example":"data"}`),
	}

	var output map[string]interface{}

	err := res.ScanJSON(&output)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, output["example"], "data")
}

func TestResource_ScanBinary(t *testing.T) {
	res := Resource{
		Data: []byte{0x3, 0x4, 0x0, 0x54},
	}

	var output int

	err := res.ScanBinary(&output)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, output, 42)
}

func TestResource_ScanBase64(t *testing.T) {
	res := Resource{
		Data: []byte("DwwADGV4YW1wbGUgZGF0YQ=="),
	}

	var output string

	err := res.ScanBase64(&output)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, output, "example data")
}
