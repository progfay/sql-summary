all: build

build:
	go build -o sqlsummary ./cmd/sqlsummary/main.go
