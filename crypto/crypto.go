package licensekey

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// GenerateKeys generates RSA public and private keys with a given bit strength.
func GenerateKeys(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	privateKey.Precompute()
	if err := privateKey.Validate(); err != nil {
		return nil, err
	}

	return privateKey, nil
}

// StorePrivateKey stores a RSA private key as PEM on disk.
func StorePrivateKey(privateKey *rsa.PrivateKey, path string) error {
	der := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}
	data := pem.EncodeToMemory(block)
	return ioutil.WriteFile(path, data, 0600)
}

// StorePublicKey stores a RSA public key as PEM on disk.
func StorePublicKey(publicKey *rsa.PublicKey, path string) error {
	der, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	block := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: der}
	data := pem.EncodeToMemory(block)
	return ioutil.WriteFile(path, data, 0600)
}

// LoadPrivateKey loads a RSA private key as PEM from disk.
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return PrivateKeyFromBytes(data)
}

// LoadPublicKey loads a RSA public key as PEM from disk.
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return PublicKeyFromBytes(data)
}

// PrivateKeyFromBytes parses a RSA private key as PEM from bytes.
func PrivateKeyFromBytes(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("no valid PEM data found")
	}

	der := block.Bytes
	privateKey, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("private key can't be decoded: %s", err)
	}

	return privateKey, nil
}

// PublicKeyFromBytes parses a RSA public key as PEM from bytes.
func PublicKeyFromBytes(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("no valid PEM data found")
	}

	der := block.Bytes
	publicKey, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("public key can't be decoded: %s", err)
	}

	return publicKey.(*rsa.PublicKey), nil
}

// Sign RSA PKCS#1 v1.5 signs a message with a given private key.
func Sign(key *rsa.PrivateKey, msg []byte) ([]byte, error) {
	hashed := sha256.Sum256(msg)
	return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
}

// Verify RSA PKCS#1 v1.5 verifies a signature against a message with a given
// public key.
func Verify(key *rsa.PublicKey, msg []byte, sig []byte) error {
	hashed := sha256.Sum256(msg)
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], sig)
}
