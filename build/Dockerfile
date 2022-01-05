FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh

RUN mkdir /etc/kira
WORKDIR /etc/kira

ADD kira /etc/kira

RUN chmod 655 /etc/kira/kira

ENTRYPOINT ["/etc/kira/kira"]
EXPOSE 8080
