all: build

clean:
	rm -rf ./build/*

run: clean build
	./build/gogit

build: clean
	GOOS=linux GOARCH=arm GOARM=5 go build -o ./build/gogit-arm5
	mkdir ./build/templates
	cp ./templates/index.html ./build/templates/index.html

install:
	go get -v .

deploy:
	scp ./build/gogit-arm5 root@192.168.0.249:/volume1/storage/opt/
	scp ./build/templates/index.html root@192.168.0.249:/volume1/storage/opt/templates/index.html
