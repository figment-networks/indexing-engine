package indexing

import (
	"github.com/figment-networks/indexing-engine/pipeline"
	"sync"
)

var (
	payloadPool = sync.Pool{
		New: func() interface{} {
			return new(payload)
		},
	}

	_ pipeline.PayloadFactory = (*payloadFactory)(nil)
	_ pipeline.Payload        = (*payload)(nil)
)

func NewPayloadFactory() *payloadFactory {
	return &payloadFactory{}
}

type payloadFactory struct{}

func (pf *payloadFactory) GetPayload(currentHeight int64) pipeline.Payload {
	payload := payloadPool.Get().(*payload)
	payload.CurrentHeight = currentHeight
	return payload
}

type payload struct {
	CurrentHeight int64
}

func (p *payload) SetCurrentHeight(height int64) {
	p.CurrentHeight = height
}

func (p *payload) GetCurrentHeight() int64 {
	return p.CurrentHeight
}

func (p *payload) MarkAsProcessed() {
	payloadPool.Put(p)
}
