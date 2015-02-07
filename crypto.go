package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
)

// AES是对称加密算法
// Key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
// 下面的AES使用CBC模式和PKCS5Padding填充法
// AES加密，传入plaintext和key，返回ciphertext（plaintext不改变）
func AesEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	ciphertext := make([]byte, aes.BlockSize+len(plaintext)+padding)

	// 初始化向量IV放在ciphertext前面
	iv := ciphertext[:aes.BlockSize]
	io.ReadFull(rand.Reader, iv)

	copy(ciphertext[aes.BlockSize:], plaintext)

	// PKCS5Padding填充
	for i := 0; i < padding; i++ {
		ciphertext[aes.BlockSize+len(plaintext)+i] = byte(padding)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}

// AES解密，传入ciphertext和key，返回plaintext（ciphertext不改变）
func AesDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	unpadding := int(plaintext[len(plaintext)-1])
	if unpadding <= 0 || unpadding > 16 {
		return nil, errors.New("AesDecrypt error: unpadding <= 0 || unpadding > 16")
	}

	plaintext = plaintext[:len(plaintext)-unpadding]

	return plaintext, nil
}

// RSA加密，传入plaintext和publickey，返回ciphertext（plaintext不改变）
func RsaEncrypt(plaintext, publickey []byte) ([]byte, error) {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, plaintext)
}

// RSA解密，传入ciphertext和privatekey，返回plaintext（ciphertext不改变）
func RsaDecrypt(ciphertext, privatekey []byte) ([]byte, error) {
	block, _ := pem.Decode(privatekey)
	if block == nil {
		return nil, errors.New("private key error")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
