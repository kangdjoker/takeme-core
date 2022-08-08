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

func RSADecrypDashboard(encryptedString string) (string, error) {

	privateKeyFileByte, err := ioutil.ReadFile("rsa_1024_priv.pem")
	if err != nil {
		return "", ErrorInternalServer(AesError, "AES Error")
	}

	b, _ := pem.Decode(privateKeyFileByte)

	privateKey, error := x509.ParsePKCS1PrivateKey(b.Bytes)
	// privateKey, error := ssh.ParseRawPrivateKey(b.Bytes)
	if error != nil {
		return "", ErrorInternalServerAesError, "AES Error")
	}

	actualPrivateKey := privateKey

	base64DecodeBytes, err := base64.StdEncoding.DecodeString(encryptedString)
	if error != nil {
		return "", model.InternalServer(model.AesError, "AES Error")
	}

	decryptedData, decryptErr := rsa.DecryptPKCS1v15(rand.Reader, actualPrivateKey, base64DecodeBytes)
	if decryptErr != nil {
		return "", model.InternalServer(model.AesError, "AES Error")
	}

	return string(decryptedData), nil
}
