#!/bin/bash

PROJECT=johnbelamaric/contacts-app
BINARY=contacts-api
cd api && docker run -v $(pwd):/go/src/github.com/$PROJECT infoblox/buildtool sh -c "cd /go/src/github.com/$PROJECT && go get && go build -o $BINARY" && \
cd .. && docker build -t johnbelamaric/contacts-api . && \
docker push johnbelamaric/contacts-api
