package datalake

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
)

// Resource represents an object being stored
type Resource struct {
	Name string
	Data []byte
}

// NewResource creates a resource
func NewResource(name string, data []byte) *Resource {
	return &Resource{
		Name: name,
		Data: data,
	}
}

// NewJSONResource creates a JSON resource
func NewJSONResource(name string, obj interface{}) (*Resource, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return &Resource{
		Name: name,
		Data: data,
	}, nil
}

// NewBinaryResource creates a binary resource
func NewBinaryResource(name string, obj interface{}) (*Resource, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	return &Resource{
		Name: name,
		Data: buf.Bytes(),
	}, nil
}

// NewBase64Resource creates a Base64 resource
func NewBase64Resource(name string, obj interface{}) (*Resource, error) {
	res, err := NewBinaryResource(name, obj)
	if err != nil {
		return nil, err
	}

	length := base64.StdEncoding.EncodedLen(len(res.Data))
	data := make([]byte, length)

	base64.StdEncoding.Encode(data, res.Data)

	return &Resource{
		Name: name,
		Data: data,
	}, nil
}

// ScanJSON parses the resource data as JSON
func (r *Resource) ScanJSON(obj interface{}) error {
	return json.Unmarshal(r.Data, obj)
}

// ScanBinary parses the resource data as binary
func (r *Resource) ScanBinary(obj interface{}) error {
	dec := gob.NewDecoder(bytes.NewReader(r.Data))

	return dec.Decode(obj)
}

// ScanBase64 parses the resource data as Base64
func (r *Resource) ScanBase64(obj interface{}) error {
	length := base64.StdEncoding.DecodedLen(len(r.Data))
	res := Resource{Data: make([]byte, length)}

	base64.StdEncoding.Decode(res.Data, r.Data)

	return res.ScanBinary(obj)
}
