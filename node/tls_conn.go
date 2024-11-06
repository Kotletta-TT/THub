package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"strings"

	"software.sslmate.com/src/go-pkcs12"
)

// Сборка TLS конфигурации
func GetTLSConfig(CAPool *x509.CertPool, cert []tls.Certificate) *tls.Config {
	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		RootCAs:            CAPool,
		InsecureSkipVerify: false,
		Certificates:       cert,
	}
}

// Чтение CA сертификата и возврат его в виде *x509.CertPool
func GetCACertPool(CAPath string) *x509.CertPool {
	caCert, err := os.ReadFile(CAPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool
}

// Deprecated: использовать GetCertKeyPairFromPFX
// Получение сертификата и ключа из файловой системы
func GetCertKeyPair(certPath, keyPath string) ([]tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	return []tls.Certificate{cert}, err
}

func ConvertX509toTLS(privateKey interface{}, cert *x509.Certificate) tls.Certificate {
	return tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privateKey,
		Leaf:        cert,
	}
}

// Получение сертификата и ключа из pfx/p12 хранилища
func GetCertKeyPairFromPFX(pfxPath, password string) ([]tls.Certificate, error) {
	pfxData, err := os.ReadFile(pfxPath)
	if err != nil {
		return []tls.Certificate{}, err
	}
	privateKey, cert, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return []tls.Certificate{}, err
	}
	tlsCert := ConvertX509toTLS(privateKey, cert)
	return []tls.Certificate{tlsCert}, nil
}

// Получение пароля для pfx/p12 хранилища
func GetKeystorePassword(configPassword string) string {
	if configPassword != "" {
		return configPassword
	}
	buf, err := os.ReadFile("/etc/edge-monitoring/ssl/password")
	if err != nil {
		log.Fatalf("keystore password file err: %s\n", err)
	}
	return strings.TrimSpace(string(buf[:]))
}

func TLSNew(pfxPath, password, ca string) (*tls.Config, error) {
	cert, err := GetCertKeyPairFromPFX(pfxPath, password)
	if err != nil {
		return nil, err
	}
	return GetTLSConfig(GetCACertPool(ca), cert), nil
}
