#!/bin/bash

KEYS_PATH="./keys"
APP_PKG="${1}"
APP_PKG_SIG="${APP_PKG}.sig"

# Verify the signature using the corresponding public key
openssl dgst -sha256 -verify "${KEYS_PATH}/public.pem" -signature "${APP_PKG_SIG}" "${APP_PKG}"
if [ $? -eq 0 ]; then
    echo "Signature is valid."
else
    echo "Signature is invalid."
fi