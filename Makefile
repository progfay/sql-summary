clean:
	rm sqlsummary

sqlsummary: clean
	go build -o sqlsummary ./cmd/sqlsummary/main.go
