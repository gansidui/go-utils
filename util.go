package utils

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"hash"
	"io"
	mrand "math/rand"
	"os"
	"reflect"
	"time"
	"unsafe"
)

// 计算string的md5值，以32位字符串形式返回
func StringToMd5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

// 计算[]byte的md5值，以32位字符串形式返回
func BytesToMd5(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// 计算文件的md5值，以32位字符串形式返回
func FileToMd5(name string) (string, error) {
	h := md5.New()
	if err := readFile(name, h); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 时间戳转换为string显示
func TimestampToString(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// 辅助函数，扫描文件内容并编码到hash.Hash中
func readFile(name string, h hash.Hash) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(name)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(file)
	for s.Scan() {
		h.Write(s.Bytes())
	}

	return s.Err()
}

// AES是对称加密算法
// Key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
// 这里使用CBC加密模式和PKCS5Padding填充法
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
	io.ReadFull(crand.Reader, iv)

	copy(ciphertext[aes.BlockSize:], plaintext)

	// PKCS5Padding填充
	for i := 0; i < padding; i++ {
		ciphertext[aes.BlockSize+len(plaintext)+i] = byte(padding)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}

// AES解密，传入ciphertext和key，返回plaintext（ciphertext改变）
func AesDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	unpadding := int(ciphertext[len(ciphertext)-1])
	ciphertext = ciphertext[:len(ciphertext)-unpadding]

	return ciphertext, nil
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

	return rsa.EncryptPKCS1v15(crand.Reader, pub, plaintext)
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

	return rsa.DecryptPKCS1v15(crand.Reader, priv, ciphertext)
}

// 得到一个长度在区间[m, n]内的随机字符串，字母为小写[a, z]
func RandomString(m, n int) string {
	num := 0
	if m < n {
		num = mrand.Intn(n-m) + m
	} else {
		num = m
	}

	bytes := make([]byte, num)
	const alphabet = "abcdefghijklmnopqrstuvwxyz"

	for i, _ := range bytes {
		bytes[i] = alphabet[mrand.Intn(26)]
	}

	return string(bytes)
}

// 不需要拷贝即可返回字符串 *s 的 byte slice，但是不能对返回的byte slice做任何修改，否则panic
func StringToByteSlice(s *string) []byte {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(s))
	sh.Cap = sh.Len
	return *(*[]byte)(unsafe.Pointer(sh))
}
