BIN_NAME = myip

build:
	go build -o ./bin/$(BIN_NAME)

.PHONY: run
run: build
	./bin/$(BIN_NAME)

test:
	go test -v ./...

clean:
	rm -f ./bin/$(BIN_NAME)
