package structs

import (
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("record not found")
)

// BlockWithMeta contains the block deails with metadata
type BlockWithMeta struct {
	// Network name
	Network string `json:"network,omitempty"`
	// ChainID
	ChainID string `json:"chain_id,omitempty"`
	// Version of block details
	Version string `json:"version,omitempty"`
	// Block details
	Block Block `json:"block,omitempty"`
}

// Block contains the block details
type Block struct {
	// ID
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt of block creation time in database
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// UpdatedAt of block update time in database
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// Hash of the Block
	Hash string `json:"hash,omitempty"`
	// Height of the Block
	Height uint64 `json:"height,omitempty"`
	// Time of the Block
	Time    time.Time `json:"time,omitempty"`
	Epoch   string    `json:"epoch,omitempty"`
	ChainID string    `json:"chain_id,omitempty"`

	NumberOfTransactions uint64 `json:"num_txs,omitempty"`
}

type TransactionSearch struct {
	Network  string   `json:"network"`
	ChainIDs []string `json:"chain_ids"`
	Epoch    string   `json:"epoch"`

	Height     uint64    `json:"height"`
	Type       []string  `json:"type"`
	BlockHash  string    `json:"block_hash"`
	Hash       string    `json:"hash"`
	Account    []string  `json:"account"`
	Sender     []string  `json:"sender"`
	Receiver   []string  `json:"receiver"`
	Memo       string    `json:"memo"`
	BeforeTime time.Time `json:"before_time"`
	AfterTime  time.Time `json:"after_time"`
	Limit      uint64    `json:"limit"`
	Offset     uint64    `json:"offset"`

	AfterHeight  uint64 `form:"after_id"`
	BeforeHeight uint64 `form:"before_id"`

	WithRaw    bool `json:"with_raw"`
	WithRawLog bool `json:"with_raw_log"`
}

type TransactionWithMeta struct {
	Network     string      `json:"network,omitempty"`
	Version     string      `json:"version,omitempty"`
	ChainID     string      `json:"chain_id,omitempty"`
	Transaction Transaction `json:"transaction,omitempty"`
}

// Transaction contains the blockchain transaction details
type Transaction struct {
	// ID of transaction assigned on database write
	ID uuid.UUID `json:"id,omitempty"`
	// Created at
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Updated at
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// Hash of the transaction
	Hash string `json:"hash,omitempty"`
	// BlockHash - hash of the block of transaction
	BlockHash string `json:"block_hash,omitempty"`
	// Height - height of the block of transaction
	Height uint64 `json:"height,omitempty"`

	Epoch string `json:"epoch,omitempty"`
	// ChainID - chain id of transacion
	ChainID string `json:"chain_id,omitempty"`
	// Time - time of transaction
	Time time.Time `json:"time,omitempty"`

	// Fee - Fees for transaction (if applies)
	Fee []TransactionAmount `json:"transaction_fee,omitempty"`
	// GasWanted
	GasWanted uint64 `json:"gas_wanted,omitempty"`
	// GasUsed
	GasUsed uint64 `json:"gas_used,omitempty"`
	// Memo - the description attached to transactions
	Memo string `json:"memo,omitempty"`

	// Version - Version of transaction record
	Version string `json:"version"`
	// Events - Transaction contents
	Events TransactionEvents `json:"events,omitempty"`

	// Raw - Raw transaction bytes
	Raw []byte `json:"raw,omitempty"`

	// RawLog - RawLog transaction's log bytes
	RawLog []byte `json:"raw_log,omitempty"`

	// HasErrors - indicates if Transaction has any errors inside
	HasErrors bool `json:"has_errors"`
}

// TransactionEvents - a set of TransactionEvent
type TransactionEvents []TransactionEvent

func (te *TransactionEvents) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &te)
}

// TransactionEvent part of transaction contents
type TransactionEvent struct {
	// ID UniqueID of event
	ID string `json:"id,omitempty"`
	// The Kind of event
	Kind string `json:"kind,omitempty"`
	// Type of transaction
	Type []string `json:"type,omitempty"`
	// Collection from where transaction came from
	Module string `json:"module,omitempty"`
	// List of sender accounts with optional amounts
	// Subcontents of event
	Sub []SubsetEvent `json:"sub,omitempty"`
}

// TransactionAmount structure holding amount information with decimal implementation (numeric * 10 ^ exp)
type TransactionAmount struct {
	// Textual representation of Amount
	Text string `json:"text,omitempty"`
	// The currency in what amount is returned (if applies)
	Currency string `json:"currency,omitempty"`

	// Numeric part of the amount
	Numeric *big.Int `json:"numeric,omitempty"`
	// Exponential part of amount obviously 0 by default
	Exp int32 `json:"exp,omitempty"`
}

// SubsetEvent - structure storing main contents of transacion
type SubsetEvent struct {
	// ID UniqueID of subsetevent
	ID string `json:"id,omitempty"`
	// Type of transaction
	Type   []string `json:"type,omitempty"`
	Action string   `json:"action,omitempty"`
	// Collection from where transaction came from
	Module string `json:"module,omitempty"`
	// List of sender accounts with optional amounts
	Sender []EventTransfer `json:"sender,omitempty"`
	// List of recipient accounts with optional amounts
	Recipient []EventTransfer `json:"recipient,omitempty"`
	// The list of all accounts that took part in the subsetevent
	Node map[string][]Account `json:"node,omitempty"`
	// Transaction nonce
	Nonce string `json:"nonce,omitempty"`
	// Completion time
	Completion *time.Time `json:"completion,omitempty"`
	// List of Amounts
	Amount map[string]TransactionAmount `json:"amount,omitempty"`
	// List of Transfers with amounts and optional recipients
	Transfers map[string][]EventTransfer `json:"transfers,omitempty"`
	// Optional error if occurred
	Error *SubsetEventError `json:"error,omitempty"`
	// Set of additional parameters attached to transaction (used as last resort)
	Additional map[string][]string `json:"additional,omitempty"`
	// SubEvents because some messages are in fact carying another messages inside
	Sub []SubsetEvent `json:"sub,omitempty"`
	// List of smart contracts details
	SmartContracts []SmartContractDataEVM `json:"smart_contracts,omitempty"`
}

// SmartContractDataEVM contains the smart contract details
type SmartContractDataEVM struct {
	// Contract method
	Input EMFunction `json:"input,omitempty"`
	// Contract events
	Output []EMLogs `json:"output,omitempty"`
	// List of Internal Transactions
	Internals []EMInternal `json:"internals,omitempty"`
}

// EMFunction contains smart contract method details
type EMFunction struct {
	// Address of the smart contract
	Address string `json:"address,omitempty"`
	// Caller account address
	Caller string `json:"caller,omitempty"`
	// From account address
	From string `json:"from,omitempty"`
	// To account address
	To string `json:"to,omitempty"`
	// Transaction value
	Value big.Int `json:"value,omitempty"`
	// Undecoded input
	Input string `json:"input,omitempty"`
	// Undecoded output
	Output string `json:"output,omitempty"`

	// Execution error
	Error string `json:"error,omitempty"`
	// Method function
	Function string `json:"func,omitempty"`
	// Smart contract method function hash
	FunctionHash string `json:"funcHash,omitempty"`
	// Method name
	Name string `json:"name,omitempty"`
	// Input details
	Inputs []Argument `json:"inputs,omitempty"`
	// Output details
	Outputs []Argument `json:"outputs,omitempty"`
}

// EMLogs contains smart contract event details
type EMLogs struct {
	// Address of the smart contract
	Address string `json:"address,omitempty"`
	// Event data
	Data string `json:"data,omitempty"`
	// Event function
	Function string `json:"func,omitempty"`
	// Event name
	Name string `json:"name,omitempty"`
	// Event topics
	Topics []string `json:"topics,omitempty"`
	// Input details
	Inputs []Argument `json:"inputs,omitempty"`
	// Output details
	Outputs []Argument `json:"outputs,omitempty"`
	// Raw event
	Raw []byte `json:"raw,omitempty"`
}

type Argument struct {
	Name    string   `json:"name,omitempty"`
	Indexed bool     `json:"indexed,omitempty"`
	Order   int      `json:"order,omitempty"`
	Type    string   `json:"type,omitempty"`
	Value   []string `json:"value,omitempty"`
}

type EMInternal struct {
	*EMFunction
	Type string `json:"type,omitempty"`
	Raw  []byte `json:"raw,omitempty"`
}

// EventTransfer - Account and Amounts pair
type EventTransfer struct {
	// Account recipient
	Account Account `json:"account,omitempty"`
	// Amounts from Transfer
	Amounts []TransactionAmount `json:"amounts,omitempty"`
}

// Account - Extended Account information
type Account struct {
	// Unique account identifier
	ID string `json:"id"`
	// External optional account details (if applies)
	Details *AccountDetails `json:"detail,omitempty"`
}

// AccountDetails External optional account details (if applies)
type AccountDetails struct {
	// Description of account
	Description string `json:"description,omitempty"`
	// Contact information
	Contact string `json:"contact,omitempty"`
	// Name of account
	Name string `json:"name,omitempty"`
	// Website address
	Website string `json:"website,omitempty"`
}

// SubsetEventError error structure for event
type SubsetEventError struct {
	// Message from error event
	Message string `json:"message,omitempty"`
}

type TransactionSearchInternal struct {
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
