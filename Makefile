.PHONY: generate-mocks test

generate-mocks:
	@echo "[mockgen] generating mocks"
	@mockgen -destination pipeline/mock/mocks.go github.com/figment-networks/indexing-engine/pipeline PayloadFactory,Payload,Source,Sink,Stage,StageRunner,Task

test:
	@echo "[go test] running tests and collecting coverage metrics"
	@go test -v -tags all_tests -race -coverprofile=coverage.txt -covermode=atomic ./...