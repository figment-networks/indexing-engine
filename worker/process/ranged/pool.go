package ranged

import (
	"sync"
)

var (
	oHBTxPool = NewHBTxPool(20)
)

type hBTxPool struct {
	stor chan chan hBTx
	lock *sync.Mutex
}

func NewHBTxPool(cap int) *hBTxPool {
	return &hBTxPool{
		stor: make(chan chan hBTx, cap),
		lock: &sync.Mutex{},
	}
}

func (o *hBTxPool) Get() chan hBTx {
	o.lock.Lock()
	defer o.lock.Unlock()
	select {
	case a := <-o.stor:
		// (lukanus): better safe than sorry
		hBTxDrain(a)
		return a
	default:
	}

	return make(chan hBTx, 10)
}

func (o *hBTxPool) Put(or chan hBTx) {
	o.lock.Lock()
	defer o.lock.Unlock()
	select {
	case o.stor <- or:
	default:
		close(or)
	}
}

func hBTxDrain(c chan hBTx) {
	for {
		select {
		case <-c:
		default:
			return
		}
	}
}
