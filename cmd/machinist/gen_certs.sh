#!/bin/bash

mkdir -p ./pki

openssl req -x509 -newkey rsa:2048 -sha256 -days 3650 -nodes \
  -keyout ./pki/server.key -out ./pki/server.crt \
  -subj '/CN=test-server' \
  -extensions san \
  -config <(echo '[req]'; echo 'distinguished_name=req';
            echo '[san]'; echo 'subjectAltName=URI:urn:127.0.0.1:test-server,IP:127.0.0.1')

cat ./pki/server.crt >> ./pki/trusted.pem