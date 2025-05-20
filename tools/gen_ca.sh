#!/bin/bash
set -e

cat > ca_openssl.cnf <<EOF
[ req ]
default_bits       = 2048
prompt             = no
distinguished_name = dn
x509_extensions    = v3_ca

[ dn ]
CN = My Root CA

[ v3_ca ]
subjectKeyIdentifier   = hash
authorityKeyIdentifier = keyid:always,issuer
basicConstraints       = critical,CA:true
keyUsage               = critical,keyCertSign,cRLSign
EOF

openssl genrsa -out ca.key 2048
openssl req -x509 -new -days 3650 -key ca.key -out ca.crt -config ca_openssl.cnf

cat > server_openssl.cnf <<EOF
[ req ]
default_bits       = 2048
prompt             = no
distinguished_name = dn
req_extensions     = v3_req

[ dn ]
CN = localhost

[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = critical,digitalSignature,keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF

openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config server_openssl.cnf

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650 -extfile server_openssl.cnf -extensions v3_req

cat > client_openssl.cnf <<EOF
[ req ]
default_bits       = 2048
prompt             = no
distinguished_name = dn
req_extensions     = v3_req

[ dn ]
CN = testclient

[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = critical,digitalSignature,keyEncipherment
extendedKeyUsage = clientAuth
EOF

openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -config client_openssl.cnf

openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 3650 -extfile client_openssl.cnf -extensions v3_req

echo "Done: "
ls -l ca.* server.* client.*