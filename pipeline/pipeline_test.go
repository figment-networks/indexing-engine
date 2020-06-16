package pipeline_test

import (
	"context"
	"github.com/figment-networks/indexing-engine/pipeline"
	mock "github.com/figment-networks/indexing-engine/pipeline/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

var (
	_ pipeline.Source = (*sourceMock)(nil)
)

type payloadMock struct{}

func (p *payloadMock) MarkAsProcessed() {}

type sourceMock struct {
	startHeight   int64
	endHeight     int64
	currentHeight int64
}

func(s *sourceMock) Next(context.Context, pipeline.Payload) bool {
	if s.currentHeight == s.endHeight {
		return false
	}
	s.currentHeight = s.currentHeight + 1
	return true
}

func(s *sourceMock) Current() int64 {
	return s.currentHeight
}


func (s *sourceMock) Err() error {
	return nil
}

func TestPipeline_SetStages(t *testing.T) {
	t.Run("all set stages are executed", func(t *testing.T) {
		ctrl, ctx := gomock.WithContext(context.Background(), t)
		defer ctrl.Finish()

		payloadFactoryMock := mock.NewMockPayloadFactory(ctrl)
		payloadFactoryMock.EXPECT().GetPayload(gomock.Any()).Return(&payloadMock{}).Times(2)

		p := pipeline.New(payloadFactoryMock)

		setupTaskMock := mock.NewMockTask(ctrl)
		fetcherTaskMock := mock.NewMockTask(ctrl)
		parserTaskMock := mock.NewMockTask(ctrl)
		validatorTaskMock := mock.NewMockTask(ctrl)
		syncerTaskMock := mock.NewMockTask(ctrl)
		sequencerTaskMock := mock.NewMockTask(ctrl)
		aggregatorTaskMock := mock.NewMockTask(ctrl)
		cleanupTaskMock := mock.NewMockTask(ctrl)

		setupTaskMock.EXPECT().GetName().Return("setupTask").Times(2)
		fetcherTaskMock.EXPECT().GetName().Return("fetcherTask").Times(2)
		parserTaskMock.EXPECT().GetName().Return("parserTask").Times(2)
		validatorTaskMock.EXPECT().GetName().Return("validatorTask").Times(2)
		syncerTaskMock.EXPECT().GetName().Return("syncerTask").Times(2)
		sequencerTaskMock.EXPECT().GetName().Return("sequencerTask").Times(2)
		aggregatorTaskMock.EXPECT().GetName().Return("aggregatorTask").Times(2)
		cleanupTaskMock.EXPECT().GetName().Return("cleanupTask").Times(2)

		p.SetStage(pipeline.StageSetup, pipeline.SyncRunner(setupTaskMock))
		p.SetStage(pipeline.StageFetcher, pipeline.SyncRunner(fetcherTaskMock))
		p.SetStage(pipeline.StageParser, pipeline.SyncRunner(parserTaskMock))
		p.SetStage(pipeline.StageValidator, pipeline.SyncRunner(validatorTaskMock))
		p.SetStage(pipeline.StageSyncer, pipeline.SyncRunner(syncerTaskMock))
		p.SetStage(pipeline.StageSequencer, pipeline.SyncRunner(sequencerTaskMock))
		p.SetStage(pipeline.StageAggregator, pipeline.SyncRunner(aggregatorTaskMock))
		p.SetStage(pipeline.StageCleanup, pipeline.SyncRunner(cleanupTaskMock))

		setupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		fetcherTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		parserTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		validatorTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		syncerTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		sequencerTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		aggregatorTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		cleanupTaskMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		sinkMock := mock.NewMockSink(ctrl)
		sinkMock.EXPECT().Consume(gomock.Any(), gomock.Any()).Return(nil).Times(2)

		options := &pipeline.Options{}

		if err := p.Start(ctx, &sourceMock{1, 2, 1}, sinkMock, options); err != nil {
			t.Errorf("should not return error")
		}
	})
}
