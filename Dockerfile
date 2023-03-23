FROM ubuntu:latest
ARG TARGETARCH

ENV DEBIAN_FRONTEND=noninteractive

# Replace default ubuntu source with mirrors.
# RUN sed -i 's/htt[p|ps]:\/\/archive.ubuntu.com\/ubuntu\//mirror:\/\/mirrors.ubuntu.com\/mirrors.txt/g' /etc/apt/sources.list

RUN apt-get update --fix-missing \
    && apt-get install curl pkg-config libseccomp-dev gcc -y \
    && apt-get clean

# Golang 1.19.4 installation
RUN curl -O https://dl.google.com/go/go1.19.4.linux-${TARGETARCH}.tar.gz \
    && tar -C /usr/local -xzf go1.19.4.linux-${TARGETARCH}.tar.gz \
    && rm go1.19.4.linux-${TARGETARCH}.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# Set the working directory and copy all files to it.
WORKDIR /kira
COPY . .

# Build the blocksyscalls binary and move it to the executable binaries.
RUN gcc ./scripts/blocksyscalls.c -O2 -Wall -lseccomp -o blocksyscalls
RUN mv blocksyscalls /usr/local/bin/

# Build kira and start the rest server.
RUN go mod tidy
RUN go build -o main rest/main.go
CMD ["/kira/main"]