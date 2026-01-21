.PHONY: build install-python install clean run

BINARY_NAME=spectre

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

run: build
	./$(BINARY_NAME)
