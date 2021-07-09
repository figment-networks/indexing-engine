package store

import (
	"context"

	"github.com/figment-networks/indexing-engine/structs"
)

type StoreCaller interface {
	GetSearchSession(ctx context.Context) (SearchStore, error)
	GetRewardsSession(ctx context.Context) (RewardStore, error)
}

type RewardStore interface {
	StoreClaimedRewards(ctx context.Context, rwds []structs.ClaimedReward) error
	StoreUnclaimedRewards(ctx context.Context, rwds []structs.UnclaimedReward) error
}

type SearchStore interface {
	StoreTransactions(ctx context.Context, txs []structs.TransactionWithMeta) error
	StoreBlocks(ctx context.Context, blocks []structs.BlockWithMeta) error
	ConfirmHeights(ctx context.Context, heights []structs.BlockWithMeta) error
}
