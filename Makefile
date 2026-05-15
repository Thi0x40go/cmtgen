BINARY_NAME=cmtgen
GOPATH=$(shell go env GOPATH)

build:
	@echo "🛠️  Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) main.go

install:
	@echo "🚀 Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	go build -o ~/.local/bin/$(BINARY_NAME) main.go
	@echo "✅ Installed! Make sure ~/.local/bin/ is in your PATH."

clean:
	rm -rf bin/

.PHONY: build install clean
