

.PHONY: test coverage run


test:
	go test -v

coverage.out: *.go
	go test -coverprofile=coverage.out

coverage: coverage.out
	go tool cover -html=coverage.out

run:
	go run ./cmd/waypoint
