all: build

build:
	go build -o build/recovery 

clean:
	rm -rf build

test:
	go test ./... --count 1 --cover

.PHONY: build clean test