run:
	@go build -o bin/run main.go
	@./bin/run

pulser:
	@go build -o bin/pulser src/pkg/cmd/pulser/main.go
	@./bin/pulser

mig:
	@go build -o bin/migrate src/pkg/cmd/migrate/main.go
	@./bin/migrate db migrate

drop:
	@go build -o bin/migrate src/pkg/cmd/migrate/main.go
	@./bin/migrate db rollback

migrate:
	@go build -o bin/migrate src/pkg/cmd/migrate/main.go
	@./bin/migrate db $(cmd) $(name)
