package ccrypto

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func TestPrivateKey2PEM(t *testing.T) {
	// 生成一个新的私钥
	_, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// 调用 privateKey2PEM 函数
	privateKeyPEM := privateKey2PEM(privateKey)

	// 检查返回的 PEM 是否可以被解码
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		t.Fatalf("Failed to decode PEM")
	}

	// 检查解码后的私钥是否可以被解析为 PKCS#8 私钥
	_, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse PKCS#8 private key: %v", err)
	}
}
