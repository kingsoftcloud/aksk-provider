package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"k8s.io/client-go/util/keyutil"
)

// RsaEncrypt function to encrypt the sk using RSA
func RsaEncrypt(plainData string, publicKey *rsa.PublicKey) (string, error) {
	// RSA encryption works by encrypting using the public key
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainData))
	if err != nil {
		return "", err
	}
	// Here, we'll return the base64 encoded encrypted data
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// RsaDecrypt function to decrypt the sk using RSA
func RsaDecrypt(ciphertext string, privateKey *rsa.PrivateKey) (string, error) {
	// If the ciphertext is base64 encoded, decode it
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	// RSA decryption works by decrypting using the private key
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertextBytes)
	if err != nil {
		return "", err
	}
	// Return the decrypted data as a string
	return string(decryptedData), nil
}

func LoadPublicKey(publicKeyText string) (*rsa.PublicKey, error) {
	// If the publicKeyText is base64 encoded, decode it
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyText)
	if err != nil {
		// If decoding Base64 fails, use the original publicKeyText directly
		publicKeyBytes = []byte(publicKeyText)
	}
	pubKeys, err := keyutil.ParsePublicKeysPEM(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("data does not contain a valid RSA or ECDSA public key")
	}
	// Type assertion to ensure it's an RSA public key
	if len(pubKeys) == 0 {
		return nil, fmt.Errorf("no public keys found in the data")
	}

	p := pubKeys[0].(*rsa.PublicKey)

	return p, nil
}

func LoadPrivateKey(privateKeyText string) (*rsa.PrivateKey, error) {
	// If the privateKeyText is base64 encoded, decode it
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyText)
	if err != nil {
		// If decoding Base64 fails, use the original privateKeyText directly
		privateKeyBytes = []byte(privateKeyText)
	}
	// Parse the private key (either base64 decoded or already PEM format)
	privKey, err := keyutil.ParsePrivateKeyPEM(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("data does not contain a valid RSA or ECDSA private key")
	}

	// Allow RSA and ECDSA formats only
	key := new(rsa.PrivateKey)
	switch k := privKey.(type) {
	case *rsa.PrivateKey:
		key = k
	default:
		return nil, fmt.Errorf("the private key is not in RSA format")
	}
	return key, nil
}
