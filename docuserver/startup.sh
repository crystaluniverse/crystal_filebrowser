#!/bin/sh
export PRODUCT_NAME=documentserver
export DEBIAN_FRONTEND=noninteractive
export LANGUAGE=en_US:en
export COMPANY_NAME=onlyoffice
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
export LANG=en_US.UTF-8

chown -R root:ssl-cert /etc/ssl/private/ssl-cert-snakeoil.key
chmod 640 /etc/ssl/private/ssl-cert-snakeoil.key

echo "127.0.0.1 localhost" >> /etc/hosts
echo "::1 localhost ip6-localhost ip6-loopback" >> /etc/hosts

/app/ds/run-document-server.sh