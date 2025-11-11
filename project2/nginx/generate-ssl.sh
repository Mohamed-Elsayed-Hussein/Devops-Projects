#!/bin/sh

if [ -f  key.pem ] && [ -f cert.pem ];then
    echo "SSL certificate and key already exist. Skipping generation."
else
    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem \
        -sha256 -days 3650 -nodes \
        -subj "/C=EG/ST=Cairo/L=Cairo/O=GoApplicatio/OU=ITservice/CN=www.goapplication.com"
fi

