package main

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var deviceWhitelist = map[string]bool{
	"D2C00A6AFF09407A30D716B1083B6D8E6866FF5637BCBB4C107148BAC32A02A7": true,
}

func mTLSAuthMidware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			http.Error(w, "No certificate provided", http.StatusUnauthorized)
			return
		}
		clientCert := r.TLS.PeerCertificates[0]
		fingerprint := sha256.Sum256(clientCert.Raw)
		fingerprintHex := strings.ToUpper(hex.EncodeToString(fingerprint[:]))
		log.Printf("fingerprint: ", fingerprintHex)
		if !deviceWhitelist[fingerprintHex] {
			http.Error(w, "Unauthorized device", http.StatusForbidden)
			log.Printf("Unauthorized device")
			return
		}
		next(w, r)
	}
}

func loadTLSConfig(caPath, certPath, keyPath string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read CA: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load server cert/key: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates:	[]tls.Certificate{serverCert},
		ClientCAs:  	caCertPool,
		ClientAuth: 	tls.RequireAndVerifyClientCert,
		MinVersion: 	tls.VersionTLS12,
	}
	return tlsConfig, nil
}
