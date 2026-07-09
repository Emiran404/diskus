BINARY := diskus
DIST := dist
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

# Yerel derleme
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

# Yerel kurulum (~/go/bin)
.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" .

.PHONY: run
run:
	go run . $(ARGS)

.PHONY: test
test:
	go test ./...

.PHONY: vet
vet:
	go vet ./...

# Tüm platformlar için binary üret (dist/ altına)
.PHONY: release
release: clean
	@mkdir -p $(DIST)
	GOOS=darwin  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-darwin-arm64 .
	GOOS=darwin  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-darwin-amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-linux-arm64 .
	GOOS=linux   GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-windows-amd64.exe .
	@echo "Binary'ler $(DIST)/ altında hazır."

.PHONY: clean
clean:
	rm -rf $(BINARY) $(BINARY).exe $(DIST)
