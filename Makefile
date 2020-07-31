.PHONY: mockgen test

mockgen:
	@echo "[mockgen] generating mocks"
	@mockgen -destination pipeline/mock/mocks.go github.com/figment-networks/indexing-engine/pipeline PayloadFactory,Payload,Source,Sink,Stage,Task,Logger

test:
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverage.txt -covermode=atomic ./...