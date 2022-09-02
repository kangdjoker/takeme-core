package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func RSADecrypt(encryptedString string) (string, error) {

	privateKeyFileByte, err := ioutil.ReadFile("private.der")
	if err != nil {
		fmt.Println("Error")
	}

	privString := base64.StdEncoding.EncodeToString(privateKeyFileByte)
	privString = "-----BEGIN RSA PRIVATE KEY-----\n" + privString + "\n-----END RSA PRIVATE KEY-----"

	b, _ := pem.Decode([]byte(privString))

	privateKey, error := x509.ParsePKCS8PrivateKey(b.Bytes)
	if error != nil {
		return "", err
	}

	actualPrivateKey := privateKey.(*rsa.PrivateKey)

	base64DecodeBytes, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	decryptedData, decryptErr := rsa.DecryptOAEP(sha1.New(), rand.Reader, actualPrivateKey, base64DecodeBytes, nil)
	if decryptErr != nil {
		return "", decryptErr
	}

	return string(decryptedData), nil
}

func RSAEncrypt(encryptedString string) (string, error) {

	pubKeyFileByte, err := ioutil.ReadFile("public.der")
	if err != nil {
		fmt.Println("Error")
	}

	pubString := base64.StdEncoding.EncodeToString(pubKeyFileByte)
	pubString = "-----BEGIN RSA PRIVATE KEY-----\n" + pubString + "\n-----END RSA PRIVATE KEY-----"

	b, _ := pem.Decode([]byte(pubString))

	pubKey, error := x509.ParsePKCS8PrivateKey(b.Bytes)
	if error != nil {
		return "", err
	}

	actualPublicKey := pubKey.(*rsa.PrivateKey)

	base64DecodeBytes, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	encryptData, decryptErr := rsa.EncryptOAEP(sha1.New(), rand.Reader, &actualPublicKey.PublicKey, base64DecodeBytes, nil)
	if decryptErr != nil {
		return "", decryptErr
	}

	return string(encryptData), nil
}

func RSADecrypDashboard(encryptedString string) (string, error) {

	privateKeyFileByte, err := ioutil.ReadFile("rsa_1024_priv.pem")
	if err != nil {
		return "", ErrorInternalServer(DecryptError, "AES Error")
	}

	b, _ := pem.Decode(privateKeyFileByte)

	privateKey, error := x509.ParsePKCS1PrivateKey(b.Bytes)
	// privateKey, error := ssh.ParseRawPrivateKey(b.Bytes)
	if error != nil {
		return "", ErrorInternalServer(DecryptError, "AES Error")
	}

	actualPrivateKey := privateKey

	base64DecodeBytes, err := base64.StdEncoding.DecodeString(encryptedString)
	if error != nil {
		return "", ErrorInternalServer(DecryptError, "AES Error")
	}

	decryptedData, decryptErr := rsa.DecryptPKCS1v15(rand.Reader, actualPrivateKey, base64DecodeBytes)
	if decryptErr != nil {
		return "", ErrorInternalServer(DecryptError, "AES Error")
	}

	return string(decryptedData), nil
}
