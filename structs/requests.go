package structs

import (
	"time"
)

type HeightHash struct {
	Epoch  string
	Height uint64
	Hash   string

	ChainID string
	Network string
}

type HeightAccount struct {
	Epoch   string
	Height  uint64
	Account string

	ChainID string
	Network string
}

type GetAccountBalanceResponse struct {
	Height   uint64              `json:"height"`
	Balances []TransactionAmount `json:"balances"`
}

type BalanceSummary struct {
	Height uint64              `json:"height"`
	Time   time.Time           `json:"time"`
	Amount []TransactionAmount `json:"balances"`
}

type Validator string
type RewardsPerValidator map[Validator][]RewardAmount

type GetRewardResponse struct {
	Height  uint64              `json:"height"`
	Rewards RewardsPerValidator `json:"rewards"`
}

type RewardSummary struct {
	Start     uint64         `json:"start"`
	End       uint64         `json:"end"`
	Time      time.Time      `json:"time"`
	Validator Validator      `json:"validator"`
	Amount    []RewardAmount `json:"rewards"`
}

type GetAccountDelegationsResponse struct {
	Height      uint64       `json:"height"`
	Delegations []Delegation `json:"delegations"`
}

type Delegation struct {
	Delegator string       `json:"delegator"`
	Validator Validator    `json:"validator"`
	Shares    RewardAmount `json:"shares"`
	Balance   RewardAmount `json:"balance"`
}
type RewardAPRSummary struct {
	Height       uint64       `json:"start_height"`
	TimeBucket   time.Time    `json:"time_bucket"`
	Apr          string       `json:"apr"`
	Bonded       RewardAmount `json:"bonded"`
	TotalRewards RewardAmount `json:"total_rewards"`
	Validator    Validator    `json:"validator"`
}

type LatestDataRequest struct {
	Network string `json:"network"`
	ChainID string `json:"chain_id"`
	Version string `json:"version"`
	TaskID  string `json:"task_id"`

	LastHash   string    `json:"last_hash"`
	LastEpoch  string    `json:"last_epoch"`
	LastHeight uint64    `json:"last_height"`
	LastTime   time.Time `json:"last_time"`
	RetryCount uint64    `json:"retry_count"`
	Nonce      []byte    `json:"nonce"`

	SelfCheck bool `json:"self_check"`
}

func (ldr *LatestDataRequest) FromMapStringInterface(m map[string]interface{}) {
	if n, ok := m["network"]; ok {
		ldr.Network = n.(string)
	}
	if n, ok := m["chain_id"]; ok {
		ldr.ChainID = n.(string)
	}
	if n, ok := m["version"]; ok {
		ldr.Version = n.(string)
	}
	if n, ok := m["task_id"]; ok {
		ldr.TaskID = n.(string)
	}
	if n, ok := m["last_hash"]; ok {
		ldr.LastHash = n.(string)
	}
	if n, ok := m["last_epoch"]; ok {
		ldr.LastEpoch = n.(string)
	}
	if n, ok := m["last_height"]; ok {
		ldr.LastHeight = uint64(n.(float64))
	}
	if n, ok := m["last_time"]; ok {
		ldr.LastTime, _ = time.Parse(time.RFC3339, n.(string))
	}
	if n, ok := m["retry_count"]; ok {
		ldr.RetryCount = uint64(n.(float64))
	}
	if n, ok := m["nonce"]; ok {
		if nstr, ok := n.(string); ok {
			ldr.Nonce = []byte(nstr)
		}
	}
	if n, ok := m["self_check"]; ok {
		ldr.SelfCheck = n.(bool)
	}
}

type LatestDataResponse struct {
	LastHash   string    `json:"last_hash"`
	LastHeight uint64    `json:"last_height"`
	LastTime   time.Time `json:"last_time"`
	LastEpoch  string    `json:"last_epoch"`
	Nonce      []byte    `json:"nonce"`

	Error      []byte `json:"error"`
	Processing bool   `json:"processing"`
}

type SyncDataRequest struct {
	Network string `json:"network"`
	ChainID string `json:"chain_id"`
	Version string `json:"version"`
	TaskID  string `json:"task_id"`

	LastHeight  uint64 `json:"last_height"`
	FinalHeight uint64 `json:"final_height"`

	LastHash   string    `json:"last_hash"`
	LastEpoch  string    `json:"last_epoch"`
	LastTime   time.Time `json:"last_time"`
	RetryCount uint64    `json:"retry_count"`
	Nonce      []byte    `json:"nonce"`

	SelfCheck bool `json:"selfCheck"`
}

func (sdr *SyncDataRequest) FromMapStringInterface(m map[string]interface{}) {
	if n, ok := m["network"]; ok {
		sdr.Network = n.(string)
	}
	if n, ok := m["chain_id"]; ok {
		sdr.ChainID = n.(string)
	}
	if n, ok := m["version"]; ok {
		sdr.Version = n.(string)
	}
	if n, ok := m["task_id"]; ok {
		sdr.TaskID = n.(string)
	}
	if n, ok := m["last_hash"]; ok {
		sdr.LastHash = n.(string)
	}
	if n, ok := m["last_epoch"]; ok {
		sdr.LastEpoch = n.(string)
	}
	if n, ok := m["last_height"]; ok {
		sdr.LastHeight = uint64(n.(float64))
	}
	if n, ok := m["final_height"]; ok {
		sdr.FinalHeight = uint64(n.(float64))
	}
	if n, ok := m["last_time"]; ok {
		sdr.LastTime, _ = time.Parse(time.RFC3339, n.(string))
	}
	if n, ok := m["retry_count"]; ok {
		sdr.RetryCount = uint64(n.(float64))
	}

	if n, ok := m["nonce"]; ok {
		nonceS, ok1 := n.(string)
		if ok1 {
			sdr.Nonce = []byte(nonceS)
		}
	}
	if n, ok := m["self_check"]; ok {
		sdr.SelfCheck = n.(bool)
	}

}

type SyncDataResponse struct {
	LastHash   string    `json:"last_hash"`
	LastHeight uint64    `json:"last_height"`
	LastTime   time.Time `json:"last_time"`
	LastEpoch  string    `json:"last_epoch"`
	RetryCount uint64    `json:"retry_count"`
	Nonce      []byte    `json:"nonce"`
	Error      []byte    `json:"error"`

	Processing bool `json:"processing"`
}

type HeightRange struct {
	Epoch       string
	Hash        string
	StartHeight uint64
	EndHeight   uint64

	ChainID string
	Network string
}

type Heights struct {
	Heights    []uint64           `json:"heights"`
	ErrorAt    []uint64           `json:"error_at"`
	LatestData LatestDataResponse `json:"latest"`

	NumberOfTx      uint64 `json:"num_tx"`
	NumberOfHeights uint64 `json:"num_heights"`
}
