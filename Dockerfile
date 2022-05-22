FROM ubuntu:latest

RUN apt-get update --fix-missing \
    && DEBIAN_FRONTEND="noninteractive" apt-get install \
    curl xz-utils unzip golang python3 nodejs -y

WORKDIR /kira
COPY . .
RUN go build -o main rest/main.go
CMD ["/kira/main"]