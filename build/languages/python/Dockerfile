FROM ubuntu:latest

RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install -y wget

# Installing python
RUN apt-get install -y python3
RUN apt-get install -y python3-pip
RUN pip3 install numpy
