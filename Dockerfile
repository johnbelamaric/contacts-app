FROM alpine:latest
MAINTAINER John Belamaric <john@belamaric.com> @johnbelamaric

RUN apk --update add bind-tools curl && rm -rf /var/cache/apk/*

ADD api/contacts-api /contacts-api

EXPOSE 80 80/tcp
ENTRYPOINT ["/contacts-api"]
