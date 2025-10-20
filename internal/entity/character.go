package entity

import (
	"crypto/rsa"
	"crypto/x509"
)

type Character struct {
	CharacterId uint64
	Movie       string
	Name        string
	MovieID     uint64
	Cert        *x509.Certificate
	Key         *rsa.PrivateKey
}
