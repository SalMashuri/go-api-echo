package rsa

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

const (
	privateKeyPath = "config/rsa/app.rsa"
	publicKeyPath  = "config/rsa/app.rsa.pub"
)

var (
	// VerifyKey rsa from config
	VerifyKey *rsa.PublicKey

	signKey *rsa.PrivateKey
)

// InitPublicKey return *rsa.PublicKey
func InitPublicKey() (*rsa.PublicKey, error) {
	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}
	return VerifyKey, nil
}

// InitPrivateKey return *rsa.PrivateKey
func InitPrivateKey() (*rsa.PrivateKey, error) {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}
	return signKey, nil
}
