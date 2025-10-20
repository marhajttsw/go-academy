package entity

import (
	"crypto/rsa"
	"crypto/x509"
)

type Movie struct {
	ID    uint64
	Title string
	Year  int
	Cert  *x509.Certificate
	Key   *rsa.PrivateKey
}
