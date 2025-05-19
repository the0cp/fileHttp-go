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
	"00:00:00:00": true,
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

func loadTLSConfig(caPath string) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read CA: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}
	return tlsConfig, nil
}
