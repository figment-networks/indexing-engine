package worker

import (
	"context"
	"sync"

	"github.com/figment-networks/indexing-engine/structs"
)

type hBTx struct {
	Height uint64
	Last   bool
}

type OutH struct {
	Height uint64
	Block  structs.BlockWithMeta
	Error  error
}

type BTX interface {
	BlockAndTx(ctx context.Context, height uint64) (blockWM structs.BlockWithMeta, txsWM []structs.TransactionWithMeta, err error)
}

type RangeRequester struct {
	BTX     BTX
	workers int
}

func NewRangeRequester(btx BTX, workers int) *RangeRequester {
	return &RangeRequester{BTX: btx, workers: workers}
}

// getRange gets given range of blocks and transactions
func (rr *RangeRequester) GetRange(ctx context.Context, hr structs.HeightRange) (h structs.Heights, err error) {
	errored := make(chan struct{})
	fin := make(chan struct{}, rr.workers)

	outH := make(chan OutH, rr.workers)
	chH := make(chan hBTx, 10) //oHBTxPool.Get()

	lock := &sync.Mutex{}
	isErr := false

	wg := &sync.WaitGroup{}
	for i := 0; i < rr.workers; i++ {
		wg.Add(1)
		go rr.asyncBlockAndTx(ctx, wg, chH, outH, errored, lock, &isErr, fin)
	}
	go rr.populateRange(chH, hr, errored)

	outHeight := &structs.Heights{}
	var count int
RANGE_LOOP:
	for {
		select {
		case <-fin:
			count++
			if count == rr.workers {
				l := len(outH)
				if l > 0 {
				DRAIN:
					for h := range outH {
						l--
						// DRAIN
						if h.Error != nil {
							err = h.Error
							outHeight.ErrorAt = append(outHeight.ErrorAt, h.Height)
							if l == 0 {
								break DRAIN
							}
							continue
						}
						assign(outHeight, h)
						if l == 0 {
							break DRAIN
						}
					}
				}
				break RANGE_LOOP
			}
		case h, ok := <-outH:
			if !ok {
				break RANGE_LOOP
			}
			if h.Error != nil {
				err = h.Error
				outHeight.ErrorAt = append(outHeight.ErrorAt, h.Height)
				continue
			}
			assign(outHeight, h)
		case <-ctx.Done():
			break RANGE_LOOP
		}
	}
	wg.Wait()
	close(outH)
	if !isErr {
		close(errored)
	}

	return *outHeight, err
}

func assign(outHeight *structs.Heights, h OutH) {
	outHeight.Heights = append(outHeight.Heights, h.Height)
	outHeight.NumberOfHeights++
	outHeight.NumberOfTx += h.Block.Block.NumberOfTransactions

	if outHeight.LatestData.LastTime.IsZero() || outHeight.LatestData.LastHeight <= h.Height {
		outHeight.LatestData.LastEpoch = h.Block.Block.Epoch
		outHeight.LatestData.LastHash = h.Block.Block.Hash
		outHeight.LatestData.LastHeight = h.Height
		outHeight.LatestData.LastTime = h.Block.Block.Time
	}
}

func (rr *RangeRequester) asyncBlockAndTx(ctx context.Context, wg *sync.WaitGroup, cinn chan hBTx, out chan OutH, er chan struct{}, l *sync.Mutex, isErr *bool, fin chan struct{}) {
	defer wg.Done()
	for in := range cinn {
		if in.Last {
			fin <- struct{}{}
			return
		}

		b, _, err := rr.BTX.BlockAndTx(ctx, in.Height)
		l.Lock() // (lukanus): this lock is for errors from other asyncBlockAndTx
		if !*isErr {
			select {
			case _, ok := <-er:
				if !ok {
					l.Unlock()
					fin <- struct{}{}
					return
				}
			case out <- OutH{Height: in.Height, Block: b, Error: err}:
			}
			if err != nil {
				*isErr = true
				close(er)
			}
		}
		l.Unlock()
	}
	fin <- struct{}{}
}

func (rr *RangeRequester) populateRange(out chan hBTx, hr structs.HeightRange, er chan struct{}) {
	height := hr.StartHeight
POPULATE:
	for {
		select {
		case out <- hBTx{Height: height}:
		case <-er:
			break POPULATE
		}

		height++
		if height > hr.EndHeight {
			break POPULATE
		}
	}

	out <- hBTx{Last: true}

	close(out)
}
