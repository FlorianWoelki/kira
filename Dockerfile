FROM ubuntu:latest
ARG TARGETARCH

RUN apt-get update --fix-missing \
    && DEBIAN_FRONTEND="noninteractive" apt-get install \
    curl xz-utils unzip -y \
    && apt-get clean

# Golang 1.19.4 installation
RUN curl -O https://dl.google.com/go/go1.19.4.linux-${TARGETARCH}.tar.gz \
    && tar -C /usr/local -xzf go1.19.4.linux-${TARGETARCH}.tar.gz \
    && rm go1.19.4.linux-${TARGETARCH}.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /kira
COPY . .
RUN go mod tidy
RUN go build -o main rest/main.go
CMD ["/kira/main"]