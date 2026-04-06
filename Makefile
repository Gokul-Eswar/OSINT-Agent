.PHONY: build install-python install clean run test vet lint cover cover-check startup-bench check

BINARY_NAME=spectre
COVER_PROFILE=coverage.out
COVER_MIN=40

build:
	go build -o $(BINARY_NAME) cmd/spectre/main.go

install-python:
	@echo "Installing Python dependencies..."
	@if [ -f analyzer/requirements.txt ]; then \
		pip install -r analyzer/requirements.txt; \
	else \
		echo "Warning: analyzer/requirements.txt not found. Skipping Python setup."; \
	fi

install: build install-python

clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -f $(COVER_PROFILE)

run: build
	./$(BINARY_NAME)

test:
	go test ./...

vet:
	go vet ./...

lint: vet

cover:
	go test -coverprofile=$(COVER_PROFILE) ./...
	go tool cover -func=$(COVER_PROFILE)

cover-check:
	go test -coverprofile=$(COVER_PROFILE) ./...
	go run ./scripts/check_coverage.go -profile $(COVER_PROFILE) -min $(COVER_MIN)

startup-bench: build
	go run ./scripts/startup_perf -binary ./$(BINARY_NAME) -runs 20

check: lint test cover-check
