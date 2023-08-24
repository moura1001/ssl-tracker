run:
	@go build -o bin/run main.go
	@./bin/run

pulser:
	@go build -o bin/pulser src/pkg/cmd/pulser/main.go
	@./bin/pulser
