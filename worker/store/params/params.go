package params

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("record not found")
)

type TransactionSearch struct {
	Height    uint64    `json:"height"`
	Type      SearchArr `json:"type"`
	BlockHash string    `json:"block_hash"`
	Hash      string    `json:"hash"`
	Account   []string  `json:"account"`
	Sender    []string  `json:"sender"`
	Receiver  []string  `json:"receiver"`
	Memo      string    `json:"memo"`

	AfterTime  time.Time `json:"before_time"`
	BeforeTime time.Time `json:"after_time"`

	AfterHeight  uint64 `json:"after_height"`
	BeforeHeight uint64 `json:"before_height"`
	Limit        uint64 `json:"limit"`
	Offset       uint64 `json:"offset"`

	Network  string   `json:"network"`
	ChainIDs []string `json:"chain_ids"`
	Epoch    string   `json:"epoch"`

	WithRaw    bool `json:"with_raw"`
	WithRawLog bool `json:"with_raw_log"`

	HasErrors bool `json:"has_errors"`
}

type SearchArr struct {
	Value []string `json:"value"`
	Any   bool     `json:"any"`
}
