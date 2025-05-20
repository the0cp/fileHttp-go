@echo off

echo [ req ]>ca_openssl.cnf
echo default_bits = 2048>>ca_openssl.cnf
echo prompt = no>>ca_openssl.cnf
echo distinguished_name = dn>>ca_openssl.cnf
echo x509_extensions = v3_ca>>ca_openssl.cnf
echo.>>ca_openssl.cnf
echo [ dn ]>>ca_openssl.cnf
echo CN = My Root CA>>ca_openssl.cnf
echo.>>ca_openssl.cnf
echo [ v3_ca ]>>ca_openssl.cnf
echo subjectKeyIdentifier = hash>>ca_openssl.cnf
echo authorityKeyIdentifier = keyid^:always,issuer>>ca_openssl.cnf
echo basicConstraints = critical,CA^:true>>ca_openssl.cnf
echo keyUsage = critical,keyCertSign,cRLSign>>ca_openssl.cnf

openssl genrsa -out ca.key 2048
openssl req -x509 -new -days 3650 -key ca.key -out ca.crt -config ca_openssl.cnf

echo [ req ]>server_openssl.cnf
echo default_bits = 2048>>server_openssl.cnf
echo prompt = no>>server_openssl.cnf
echo distinguished_name = dn>>server_openssl.cnf
echo req_extensions = v3_req>>server_openssl.cnf
echo.>>server_openssl.cnf
echo [ dn ]>>server_openssl.cnf
echo CN = localhost>>server_openssl.cnf
echo.>>server_openssl.cnf
echo [ v3_req ]>>server_openssl.cnf
echo basicConstraints = CA^:FALSE>>server_openssl.cnf
echo keyUsage = critical,digitalSignature,keyEncipherment>>server_openssl.cnf
echo extendedKeyUsage = serverAuth>>server_openssl.cnf
echo subjectAltName = @alt_names>>server_openssl.cnf
echo.>>server_openssl.cnf
echo [ alt_names ]>>server_openssl.cnf
echo DNS.1 = localhost>>server_openssl.cnf
echo IP.1 = 127.0.0.1>>server_openssl.cnf

openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config server_openssl.cnf

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650 -extfile server_openssl.cnf -extensions v3_req

echo [ req ]>client_openssl.cnf
echo default_bits = 2048>>client_openssl.cnf
echo prompt = no>>client_openssl.cnf
echo distinguished_name = dn>>client_openssl.cnf
echo req_extensions = v3_req>>client_openssl.cnf
echo.>>client_openssl.cnf
echo [ dn ]>>client_openssl.cnf
echo CN = testclient>>client_openssl.cnf
echo.>>client_openssl.cnf
echo [ v3_req ]>>client_openssl.cnf
echo basicConstraints = CA^:FALSE>>client_openssl.cnf
echo keyUsage = critical,digitalSignature,keyEncipherment>>client_openssl.cnf
echo extendedKeyUsage = clientAuth>>client_openssl.cnf

openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -config client_openssl.cnf

openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 3650 -extfile client_openssl.cnf -extensions v3_req

echo Done:
dir ca.* server.* client.*
pause