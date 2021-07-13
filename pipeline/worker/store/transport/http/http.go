package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/figment-networks/indexing-engine/structs"
	"github.com/figment-networks/indexing-engine/worker/store"
)

type HTTPStore struct {
	cli  *http.Client
	urls []string
	next int
	len  int
	lock sync.Mutex
}

func NewHTTPStore(urls []string, cli *http.Client) *HTTPStore {
	return &HTTPStore{
		cli:  cli,
		urls: urls,
		len:  len(urls),
	}
}

func (s *HTTPStore) inc() {
	s.lock.Lock()
	if s.next == s.len-1 {
		s.next = 0
	} else {
		s.next++
	}
	s.lock.Unlock()
}

func (s *HTTPStore) GetRewardsSession(ctx context.Context) (store.RewardStore, error) {
	s.inc()
	return &RewardStore{&HTTPStoreSession{
		cli: s.cli,
		url: s.urls[s.next],
	}}, nil
}

func (s *HTTPStore) GetSearchSession(ctx context.Context) (store.SearchStore, error) {
	s.inc()
	return &SearchStore{&HTTPStoreSession{
		cli: s.cli,
		url: s.urls[s.next],
	}}, nil
}

type HTTPStoreSession struct {
	cli *http.Client
	url string
}

func (s *HTTPStoreSession) call(ctx context.Context, name string, in interface{}) error {
	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	buff.WriteString(`{"jsonrpc": "2.0", "id": 1, "method": "` + name + `", "params": `)
	if err := enc.Encode(in); err != nil {
		return err
	}
	buff.WriteString(`}`)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, buff)
	if err != nil {
		return err
	}

	resp, err := s.cli.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error calling store: returned status: %s", resp.Status)
	}

	jResp := &JsonRPCResponse{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(jResp); err != nil {
		return err
	}

	if jResp.Error != nil && jResp.Error.Message != "" {
		return fmt.Errorf("error calling store: %s", jResp.Error.Message)
	}

	return nil
}

type JsonRPCResponse struct {
	ID      uint64          `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Error   *JsonRPCError   `json:"error,omitempty"`
	Result  json.RawMessage `json:"result"`
}

type JsonRPCError struct {
	Code    int64         `json:"code"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

type SearchStore struct {
	*HTTPStoreSession
}

func (s *SearchStore) StoreTransactions(ctx context.Context, txs []structs.TransactionWithMeta) error {
	return s.call(ctx, "store_transactions", txs)
}

func (s *SearchStore) StoreBlocks(ctx context.Context, blocks []structs.BlockWithMeta) error {
	return s.call(ctx, "store_blocks", blocks)
}

func (s *SearchStore) ConfirmHeights(ctx context.Context, heights []structs.BlockWithMeta) error {
	return s.call(ctx, "confirm_heights", heights)
}

type RewardStore struct {
	*HTTPStoreSession
}

func (s *RewardStore) StoreClaimedRewards(ctx context.Context, rewards []structs.ClaimedReward) error {
	return s.call(ctx, "store_claimed_rewards", rewards)
}

func (s *RewardStore) StoreUnclaimedRewards(ctx context.Context, rewards []structs.UnclaimedReward) error {
	return s.call(ctx, "store_unclaimed_rewards", rewards)
}
