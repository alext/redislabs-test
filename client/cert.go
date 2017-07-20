package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

const keyBits = 2048

func generateCert(cn string) ([]byte, []byte, error) {

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(100 * 8760 * time.Hour), // 100 years from now
	}

	priv, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return nil, nil, fmt.Errorf("Error generating key : %s", err.Error())
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating certificate: %s", err.Error())
	}
	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return cert, key, nil
}
