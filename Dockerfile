FROM ubuntu:latest

RUN apt-get update --fix-missing \
    && DEBIAN_FRONTEND="noninteractive" apt-get install \
    curl xz-utils unzip golang python3 nodejs default-jre-headless -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /kira
COPY . .
RUN go mod tidy
RUN go build -o main rest/main.go
CMD ["/kira/main"]