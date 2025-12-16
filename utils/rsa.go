package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// 生成公钥私钥
func GenerateKey() (private, public string, err error) {
	// 生成2048位RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}

	// 导出公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// 导出私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM), string(publicKeyPEM), nil
}

func EncryptOAEP(publicKeyPEM []byte, message string) ([]byte, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub.(*rsa.PublicKey),
		[]byte(message),
		nil,
	)
}

func DecryptOAEP(privateKeyPEM []byte, ciphertext []byte) (string, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", errors.New("failed to parse PEM block")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	plaintext, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		priv,
		ciphertext,
		nil,
	)
	return string(plaintext), err
}

func SignPKCS1v15(privateKeyPEM []byte, message string) ([]byte, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256([]byte(message))
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed[:])
}

func VerifyPKCS1v15(publicKeyPEM []byte, message string, signature []byte) bool {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return false
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false
	}

	hashed := sha256.Sum256([]byte(message))
	return rsa.VerifyPKCS1v15(
		pub.(*rsa.PublicKey),
		crypto.SHA256,
		hashed[:],
		signature,
	) == nil
}
