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

// roundRobin can be used to get the next url from a list of given urls,
// starting with the first.
type roundRobin struct {
	urls []string
	next int
	len  int
	lock sync.Mutex
}

func newRoundRobin(urls []string) *roundRobin {
	return &roundRobin{
		urls: urls,
		len:  len(urls),
	}
}

// getNext returns the next url in the list.
func (r *roundRobin) getNext() string {
	r.lock.Lock()

	// get the current url
	url := r.urls[r.next]

	// increase the index
	if r.next < r.len-1 {
		r.next++
	} else {
		r.next = 0
	}

	r.lock.Unlock()

	return url
}

type HTTPStore struct {
	cli         *http.Client
	searchUrls  *roundRobin
	rewardsUrls *roundRobin
}

// NewHTTPStore constructs a new HTTPStore with the given search and rewards urls.
func NewHTTPStore(searchUrls []string, rewardsUrls []string, cli *http.Client) *HTTPStore {
	return &HTTPStore{
		cli:         cli,
		searchUrls:  newRoundRobin(searchUrls),
		rewardsUrls: newRoundRobin(rewardsUrls),
	}
}

// NewHTTPStoreSearch creates a new HTTPStore only with search urls.
func NewHTTPStoreSearch(searchUrls []string, cli *http.Client) *HTTPStore {
	return NewHTTPStore(searchUrls, nil, cli)
}

// NewHTTPStoreRewards creates a new HTTPStore only with rewards urls.
func NewHTTPStoreRewards(rewardsUrls []string, cli *http.Client) *HTTPStore {
	return NewHTTPStore(nil, rewardsUrls, cli)
}

func (s *HTTPStore) GetRewardsSession(ctx context.Context) (store.RewardStore, error) {
	return &RewardStore{&HTTPStoreSession{
		cli: s.cli,
		url: s.rewardsUrls.getNext(),
	}}, nil
}

func (s *HTTPStore) GetSearchSession(ctx context.Context) (store.SearchStore, error) {
	return &SearchStore{&HTTPStoreSession{
		cli: s.cli,
		url: s.searchUrls.getNext(),
	}}, nil
}

// HTTPStoreSession contains common logic for Store structs (SearchStore and RewardStore).
type HTTPStoreSession struct {
	cli *http.Client
	url string
}

// call sends the JSON RPC request with the given `in` data.
func (s *HTTPStoreSession) call(ctx context.Context, name string, in interface{}) error {
	// create the JSON request data
	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	buff.WriteString(`{"jsonrpc":"2.0","id":1,"method":"` + name + `","params":`)
	if err := enc.Encode(in); err != nil {
		return err
	}
	buff.WriteString(`}`)

	// send the request
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
