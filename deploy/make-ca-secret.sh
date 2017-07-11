#!/bin/bash

# generates manifest for secrets

function make_secret() {
	FILE=$1
	if [ ! -f $FILE ]; then
  		echo "$FILE not found, skipping"
	fi
	SECRET=$(basename $FILE)
  	echo "  $SECRET: $(cat $FILE | base64)"
}

if [ -z "$1" ]; then
	echo "Usage: $0 <file> [ ... ]"
	exit 1
fi

cat <<EOF
apiVersion: v1     
kind: Secret
metadata:
  name: ca
  namespace: contacts
data:
EOF

while [ ! -z "$1" ]; do
	make_secret $1
	shift;
done

cat <<EOF
type: Opaque
---
EOF
