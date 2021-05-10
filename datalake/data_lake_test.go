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

	mockStorage.EXPECT().Store(res.Data, "oasis", "mainnet", "validators")

	dl.StoreResource(res, "validators")
}

func TestDataLake_StoreResourceAtHeight(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := mock.NewMockStorage(mockCtrl)

	dl := NewDataLake("oasis", "mainnet", mockStorage)
	res := NewResource([]byte("example data"))

	mockStorage.EXPECT().Store(res.Data, "oasis", "mainnet", "height", "2000", "validators")

	dl.StoreResourceAtHeight(res, "validators", 2000)
}
