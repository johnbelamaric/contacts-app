#!/bin/bash

# generates manifest for secrets

function make_secret() {
	SECRET=$1
	FILE=$2
	if [ ! -f $FILE ]; then
  		echo "$FILE not found, skipping"
	fi
  	echo "  $SECRET: $(cat $FILE | base64)"
}

if [ -z "$1" ]; then
	echo "Usage: $0 <name> <key-file> <cert-file>"
	exit 1
fi

NAME=$1
KEY=$2
CERT=$3

shift
cat <<EOF
apiVersion: v1     
kind: Secret
metadata:
  name: $NAME
  namespace: contacts
type: kubernetes.io/tls
data:
EOF

make_secret "tls.key" $KEY
make_secret "tls.cert" $CERT
