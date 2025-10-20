package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"
)

type CAConfig struct {
	CertPath string
	KeyPath  string
}

func LoadCAFromEnv() CAConfig {
	return CAConfig{
		CertPath: os.Getenv("CA_CERT_PATH"),
		KeyPath:  os.Getenv("CA_KEY_PATH"),
	}
}

func loadPEM(certPath, keyPath string) ([]byte, []byte, error) {
	if certPath == "" || keyPath == "" {
		return nil, nil, errors.New("missing CA paths")
	}
	crt, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, err
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, err
	}
	return crt, key, nil
}

// Generate a leaf certificate signed by a given parent (cert,key in PEM).
// If parentPEM is nil, the cert will be self-signed.
func GenerateSignedCert(commonName string, parentCertPEM, parentKeyPEM []byte) (certPEM, keyPEM []byte, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return nil, nil, err
	}

	tpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
	}

	var parent *x509.Certificate
	var signer interface{}

	if len(parentCertPEM) > 0 && len(parentKeyPEM) > 0 {
		block, _ := pem.Decode(parentCertPEM)
		if block == nil {
			return nil, nil, errors.New("invalid parent cert")
		}
		pc, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, nil, err
		}
		kblock, _ := pem.Decode(parentKeyPEM)
		if kblock == nil {
			return nil, nil, errors.New("invalid parent key")
		}
		pk, err := x509.ParsePKCS1PrivateKey(kblock.Bytes)
		if err != nil {
			return nil, nil, err
		}
		parent = pc
		signer = pk
	} else {
		// self-signed
		parent = tpl
		signer = priv
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, parent, &priv.PublicKey, signer)
	if err != nil {
		return nil, nil, err
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return certPEM, keyPEM, nil
}

func GenerateMovieCert(commonName string, ca CAConfig) (certPEM, keyPEM []byte, err error) {
	crt, key, err := loadPEM(ca.CertPath, ca.KeyPath)
	if err != nil {
		return nil, nil, err
	}
	return GenerateSignedCert(commonName, crt, key)
}

func GenerateCharacterCert(commonName string, movieCertPEM, movieKeyPEM []byte) (certPEM, keyPEM []byte, err error) {
	return GenerateSignedCert(commonName, movieCertPEM, movieKeyPEM)
}

func LoadCAObjectsFromEnv() (*x509.Certificate, *rsa.PrivateKey, error) {
	ca := LoadCAFromEnv()
	certPEM, keyPEM, err := loadPEM(ca.CertPath, ca.KeyPath)
	if err != nil {
		return nil, nil, err
	}
	cblock, _ := pem.Decode(certPEM)
	if cblock == nil {
		return nil, nil, errors.New("invalid CA cert")
	}
	cert, err := x509.ParseCertificate(cblock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	kblock, _ := pem.Decode(keyPEM)
	if kblock == nil {
		return nil, nil, errors.New("invalid CA key")
	}
	key, err := x509.ParsePKCS1PrivateKey(kblock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

func GenerateSignedCertObjects(commonName string, parentCert *x509.Certificate, parentKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return nil, nil, err
	}
	tpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
	}
	if parentCert == nil || parentKey == nil {
		parentCert = tpl
		parentKey = priv
	}
	der, err := x509.CreateCertificate(rand.Reader, tpl, parentCert, &priv.PublicKey, parentKey)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, nil, err
	}
	return cert, priv, nil
}
