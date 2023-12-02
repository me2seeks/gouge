package ccrypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

func GeneratePEM() ([]byte, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return privateKey2PEM(privateKey), nil
}

func privateKey2PEM(privateKey ed25519.PrivateKey) []byte {
	// 将私钥编码为 PKCS#8 格式
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	// 将私钥编码为 PEM 格式
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return privateKeyPEM
}

func PEM2PrivateKey(privateKeyPEM []byte) (ed25519.PrivateKey, error) {
	// 将私钥从 PEM 格式解码
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		panic("failed to decode PEM block containing private key")
	}

	// 将私钥解码为 PKCS#8 格式
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return privateKey.(ed25519.PrivateKey), nil
}

// FingerprintKey calculates the SHA256 hash of an SSH public key
func FingerPrint(k ssh.PublicKey) string {
	bytes := sha256.Sum256(k.Marshal())
	return base64.StdEncoding.EncodeToString(bytes[:])
}
