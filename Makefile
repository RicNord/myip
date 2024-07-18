BIN_NAME = myip
DEFAULT_BIN_DIR := $(HOME)/go/bin
OS := $(shell uname -s)

# Check if GOPATH is set, if not use the default
ifdef GOPATH
    BIN_DIR := $(GOPATH)/bin
else
    BIN_DIR := $(DEFAULT_BIN_DIR)
endif

build:
	go build -o ./bin/$(BIN_NAME)

.PHONY: run
run: build
	./bin/$(BIN_NAME)

check:
	go fmt ./...
	go vet ./...
	go test ./...

test:
	go test -v ./...

clean:
	go clean
	rm -f ./bin/$(BIN_NAME)

install:
	go install
ifeq ($(OS),Linux)
	@mkdir -p "$(HOME)/.local/share/systemd/user"
	cp systemd/myip.service "$(HOME)/.local/share/systemd/user"
endif

uninstall:
	rm -f "$(BIN_DIR)/myip"
ifeq ($(OS),Linux)
	rm -f "$(HOME)/.local/share/systemd/user/myip.service"
endif
