#!/bin/bash
docker pull python:3.9.1-alpine
docker pull golang:1.17-alpine
docker pull gcc:latest
docker pull openjdk:8u232-jdk
docker-compose up -d
