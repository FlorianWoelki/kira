create-image:
	docker build build/all-in-one-ubuntu -t all-in-one-ubuntu

build:
	go build -o kira main.go

deploy:
	cp kira /usr/local/bin
