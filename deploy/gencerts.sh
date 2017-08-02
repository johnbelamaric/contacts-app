#!/bin/bash

#openssl genrsa -out ca-key.pem 2048
#openssl req -x509 -new -nodes -key ca-key.pem -sha256 -days 5000 -out ca.crt -subj "/CN=ngp-primer-ca"

#openssl genrsa -out contacts-api-key.pem 2048
#openssl req -new -key contacts-api-key.pem -out contacts-api.csr -subj "/CN=contacts-api" -config contacts-api.cnf
#openssl x509 -req -in contacts-api.csr -CA ca.crt -CAkey ca-key.pem -CAcreateserial -out contacts-api.pem -days 5000 -extensions v3_req -extfile contacts-api.cnf

#openssl genrsa -out minikube-key.pem 2048
#openssl req -new -key minikube-key.pem -out minikube.csr -subj "/CN=minikube" -config minikube.cnf
#openssl x509 -req -in minikube.csr -CA ca.crt -CAkey ca-key.pem -CAcreateserial -out minikube.pem -days 5000 -extensions v3_req -extfile minikube.cnf

openssl genrsa -out oauth2-key.pem 2048
openssl req -new -key oauth2-key.pem -out oauth2.csr -subj "/CN=oauth2-proxy" -config oauth2.cnf
openssl x509 -req -in oauth2.csr -CA ca.crt -CAkey ca-key.pem -CAcreateserial -out oauth2.pem -days 5000 -extensions v3_req -extfile oauth2.cnf

