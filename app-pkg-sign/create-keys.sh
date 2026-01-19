#!/bin/bash

KEYS_PATH="./keys"

# Create a new RSA private key
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:4096 -out "${KEYS_PATH}/private.pem"

# encrypt the private key with a passphrase
openssl pkcs8 -topk8 -inform PEM -outform PEM -in "${KEYS_PATH}/private.pem" -out "${KEYS_PATH}/private_encrypted.pem" -v2 aes-256-cbc

# extract the public key from the private key
openssl rsa -pubout -in "${KEYS_PATH}/private.pem" -out "${KEYS_PATH}/public.pem"