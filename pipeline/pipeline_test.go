package pipeline_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/figment-networks/indexing-engine/pipeline"
	mock "github.com/figment-networks/indexing-engine/pipeline/mock"
	"sync"
	"testing"
)

var (
	payloadPool = sync.Pool{
		New: func() interface{} {
			return new(payload)
		},
	}
)

type payload struct {
	currentHeight int64
}

func (p *payload) SetCurrentHeight(h int64) {
	p.currentHeight = h
}

func (p *payload) GetCurrentHeight() int64 {
	return p.currentHeight
}

func (p *payload) Clone() pipeline.Payload {
	newP := payloadPool.Get().(*payload)

	return newP
}

func (p *payload) MarkAsProcessed() {
	payloadPool.Put(p)
}

func TestPipeline_SetStages(t *testing.T) {
	t.Run("all set stages are executed", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		sourceMock := mock.NewMockSource(ctrl)
		sinkMock := mock.NewMockSink(ctrl)
		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)

		sourceMock.EXPECT().Run(gomock.Any()).Return(nil)
		sourceMock.EXPECT().GetStartHeight().Return(int64(1))
		sourceMock.EXPECT().GetEndHeight().Return(int64(2))
		sinkMock.EXPECT().Run(gomock.Any()).Return(nil)
		payloadFactoryMock.EXPECT().GetPayload().Return(&payload{}).Times(2)

		p := pipeline.New(sourceMock, sinkMock, payloadFactoryMock)

		setupTaskMock := mock.NewMockTask(ctrl)
		fetcherTaskMock := mock.NewMockTask(ctrl)
		parserTaskMock := mock.NewMockTask(ctrl)
		validatorTaskMock := mock.NewMockTask(ctrl)
		syncerTaskMock := mock.NewMockTask(ctrl)
		sequencerTaskMock := mock.NewMockTask(ctrl)
		aggregatorTaskMock := mock.NewMockTask(ctrl)
		cleanupTaskMock := mock.NewMockTask(ctrl)

		p.SetSetupStage(pipeline.SyncRunner(setupTaskMock))
		p.SetFetcherStage(pipeline.SyncRunner(fetcherTaskMock))
		p.SetParserStage(pipeline.SyncRunner(parserTaskMock))
		p.SetValidatorStage(pipeline.SyncRunner(validatorTaskMock))
		p.SetSyncerStage(pipeline.SyncRunner(syncerTaskMock))
		p.SetSequencerStage(pipeline.SyncRunner(sequencerTaskMock))
		p.SetAggregatorStage(pipeline.SyncRunner(aggregatorTaskMock))
		p.SetCleanupStage(pipeline.SyncRunner(cleanupTaskMock))

		setupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		fetcherTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		parserTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		validatorTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		syncerTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		sequencerTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		aggregatorTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		cleanupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		if err := p.Start(ctx); err != nil {
			t.Errorf("should not return error")
		}
	})
}
