package datalake

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/figment-networks/indexing-engine/datalake/mock"
)

func TestNewDataLake(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)

	assert.Equal(t, "oasis", dl.network)
	assert.Equal(t, "mainnet", dl.chain)

	assert.Equal(t, mockStorage, dl.storage)
}

func TestDataLake_StoreResource(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)
	res := NewResource([]byte("example data"))

	mockStorage.
		EXPECT().
		Store(res.Data, "oasis", "mainnet", "validators")

	err := dl.StoreResource(res, "validators")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDataLake_StoreResource_BlankName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)
	res := NewResource([]byte("example data"))

	err := dl.StoreResource(res, "")

	assert.Equal(t, ErrResourceNameRequired, err)
}

func TestDataLake_IsResourceStored(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)

	mockStorage.
		EXPECT().
		IsStored("oasis", "mainnet", "validators").
		Return(true, nil)

	stored, err := dl.IsResourceStored("validators")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, stored)
}

func TestDataLake_RetrieveResource(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)
	data := []byte("example data")

	mockStorage.
		EXPECT().
		Retrieve("oasis", "mainnet", "validators").
		Return(data, nil)

	res, err := dl.RetrieveResource("validators")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, res.Data)
}

func TestDataLake_StoreResourceAtHeight(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)
	res := NewResource([]byte("example data"))

	mockStorage.
		EXPECT().
		Store(res.Data, "oasis", "mainnet", "height", "2000", "transactions")

	err := dl.StoreResourceAtHeight(res, "transactions", 2000)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDataLake_StoreResourceAtHeight_BlankName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)
	res := NewResource([]byte("example data"))

	err := dl.StoreResourceAtHeight(res, "", 2000)

	assert.Equal(t, ErrResourceNameRequired, err)
}

func TestDataLake_IsResourceStoredAtHeight(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)

	mockStorage.
		EXPECT().
		IsStored("oasis", "mainnet", "height", "2000", "transactions").
		Return(true, nil)

	stored, err := dl.IsResourceStoredAtHeight("transactions", 2000)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, stored)
}

func TestDataLake_RetrieveResourceAtHeight(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)
	data := []byte("example data")

	mockStorage.
		EXPECT().
		Retrieve("oasis", "mainnet", "height", "2000", "transactions").
		Return(data, nil)

	res, err := dl.RetrieveResourceAtHeight("transactions", 2000)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, res.Data)
}
