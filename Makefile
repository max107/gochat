all: build

clean:
	rm -rf ./bin/*

run: clean build
	./bin/app

build: clean
	go build -o ./bin/app

install:
	go get -v .
