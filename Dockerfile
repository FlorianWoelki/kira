FROM ubuntu:latest

RUN apt-get update --fix-missing \
    && DEBIAN_FRONTEND="noninteractive" apt-get install \
    curl xz-utils unzip python3 nodejs -y \
    && apt-get clean

# Java installation
RUN curl "https://download.java.net/java/GA/jdk15.0.2/0d1cfde4252546c6931946de8db48ee2/7/GPL/openjdk-15.0.2_linux-x64_bin.tar.gz" -o java.tar.gz
RUN tar xzf java.tar.gz --strip-components=1
RUN rm java.tar.gz

# Golang 1.19.1 installation
RUN curl -O "https://dl.google.com/go/go1.19.1.linux-amd64.tar.gz" \
    && tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz \
    && rm go1.19.1.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /kira
COPY . .
RUN go mod tidy
RUN go build -o main rest/main.go
CMD ["/kira/main"]