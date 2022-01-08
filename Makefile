build:
	go build -o kira main.go

deploy:
	cp kira /usr/local/bin
