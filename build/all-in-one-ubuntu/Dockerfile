FROM ubuntu:latest

RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install -y wget

# Installing golang
RUN wget https://dl.google.com/go/go1.17.5.linux-arm64.tar.gz && \
    tar -xvf go1.17.5.linux-arm64.tar.gz && \
    mv go /usr/local && \
    rm -rf go1.17.5.linux-arm64.tar.gz

# Installing python
RUN apt-get install -y python3

# Installing java
RUN apt-get install -y openjdk-8-jdk

# Installing nodejs
RUN apt-get install -y nodejs

# Installing C++
RUN apt-get install -y g++
