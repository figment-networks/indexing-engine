package structs

import (
	"math/big"
	"time"
)

// RewardQuery - A set of fields used as params for search
// swagger:model
type RewardQuery struct {
	// Network identifier to search in
	//
	// required: true
	// example: cosmos
	Network string `json:"network"`
	// ChainID to search in
	//
	// required: true
	// example: 'cosmoshub-3'
	ChainID string `json:"chain_id"`
	// AfterHeight gets all transaction bigger than given height
	// Has to be bigger than BeforeHeight
	//
	// min: 0
	AfterHeight uint64 `json:"after_height"`
	// BeforeHeight gets all transaction lower than given height
	// Has to be lesser than AfterHeight
	//
	// min: 0
	BeforeHeight uint64 `json:"before_height"`
	// Account - the account identifier to look for
	Account string `json:"account"`
	// The time of transaction (if not given by chain API, the same as block)
	AfterTime time.Time `json:"after_time"`
	// The time of transaction (if not given by chain API, the same as block)
	BeforeTime time.Time `json:"before_time"`
	// Limit of how many requests to get in one request
	//
	// default: 100
	// max: 1000
	Limit uint64 `json:"limit"`
	// Offset the offset number or
	Offset uint64 `json:"offset"`
}

type ClaimedReward struct {
	// ID UniqueID of reward
	ID string `json:"id,omitempty"`
	// Created at
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Updated at
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// Reward recipient account address
	Account string `json:"account,omitempty"`
	// Chain id of reward
	ChainID string `json:"chain_id,omitempty"`
	// Validator rewards
	ClaimedReward []RewardAmount `json:"claimed_reward,omitempty"`
	// Epoch
	Epoch string `json:"epoch,omitempty"`
	// Height from which reward come from
	Height uint64 `json:"height,omitempty"`
	// Name of the network
	Network string `json:"network,omitempty"`
	// Address of validator
	Validator string `json:"validator,omitempty"`
	// Version of reward
	Version string `json:"version,omitempty"`
	// Reward time
	Time time.Time `json:"time,omitempty"`
}

type UnclaimedReward struct {
	// ID UniqueID of reward
	ID string `json:"id,omitempty"`
	// Created at
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Updated at
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// Reward recipient account address
	Account string `json:"account,omitempty"`
	// Chain id of reward
	ChainID string `json:"chain_id,omitempty"`
	// Epoch
	Epoch string `json:"epoch,omitempty"`
	// Height from which reward come from
	Height uint64 `json:"height,omitempty"`
	// Name of the network
	Network string `json:"network,omitempty"`
	// Validator rewards
	UnclaimedReward []RewardAmount `json:"unclaimed_reward,omitempty"`
	// Address of validator
	Validator string `json:"validator,omitempty"`
	// Version of reward
	Version string `json:"version,omitempty"`
	// Reward time
	Time time.Time `json:"time,omitempty"`
}

// RewardAmount structure holding amount information with decimal implementation (numeric * 10 ^ exp)
type RewardAmount struct {
	// Textual representation of Amount
	Text string `json:"text,omitempty"`
	// The currency in what amount is returned (if applies)
	Currency string `json:"currency,omitempty"`

	// Numeric part of the amount
	Numeric *big.Int `json:"numeric,omitempty"`
	// Exponential part of amount obviously 0 by default
	Exp int32 `json:"exp,omitempty"`
}

type GetClaimedRewardResponse struct {
	Height         uint64          `json:"height"`
	ClaimedRewards []ClaimedReward `json:"claimed_rewards"`
}

type GetUnclaimedRewardResponse struct {
	Height           uint64            `json:"height"`
	UnclaimedRewards []UnclaimedReward `json:"unclaimed_rewards"`
}
