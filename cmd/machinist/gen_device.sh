#!/bin/bash

mkdir -p ./pki

openssl ecparam -genkey -name prime256v1 -noout -out .\pki\device.key
openssl req -x509 -new -days 3650 -subj "/CN=unused" -key .\pki\device.key -out .\pki\device.crt 