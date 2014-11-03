package utils

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"hash"
	"io"
	"math/rand"
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
// AES加密，传入的plaintext会被重写为ciphertext，plaintext不可再利用
func AesEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	plaintext = PKCS5Padding(plaintext, blockSize)

	mode := cipher.NewCBCEncrypter(block, key[:blockSize])
	mode.CryptBlocks(plaintext, plaintext)

	return plaintext, nil
}

// AES解密，传入的ciphertext会被重写为plaintext，plaintext不可再利用
func AesDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	mode := cipher.NewCBCDecrypter(block, key[:blockSize])
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return ciphertext, nil
}

// PKCS5Padding填充法
func PKCS5Padding(b []byte, blockSize int) []byte {
	padding := blockSize - len(b)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(b, padtext...)
}

// PKCS5Padding反填充（去掉最后一个字节 unpadding 次）
func PKCS5UnPadding(b []byte) []byte {
	length := len(b)
	unpadding := int(b[length-1])
	return b[:(length - unpadding)]
}

// 得到一个长度在区间[m, n]内的随机字符串，字母为小写[a, z]
func RandomString(m, n int) string {
	rand.Seed(time.Now().UnixNano())
	num := 0
	if m < n {
		num = rand.Intn(n-m) + m
	} else {
		num = m
	}

	bytes := make([]byte, num)
	const alphabet = "abcdefghijklmnopqrstuvwxyz"

	for i, _ := range bytes {
		bytes[i] = alphabet[rand.Intn(26)]
	}

	return string(bytes)
}

// BigEndian: uint32 --> []byte
func Uint32ToBytes(v uint32) []byte {
	var b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}

// BigEndian: int32 --> []byte
func Int32ToBytes(v int32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
	return b
}

// BigEndian: []byte --> uint32
func BytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// BigEndian: []byte -->int32
func BytesToInt32(b []byte) int32 {
	return int32(b[3]) | int32(b[2])<<8 | int32(b[1])<<16 | int32(b[0])<<24
}

// 不需要拷贝即可返回字符串 *s 的 byte slice，但是不能对返回的byte slice做任何修改，否则panic
func StringToByteSlice(s *string) []byte {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(s))
	sh.Cap = sh.Len
	return *(*[]byte)(unsafe.Pointer(sh))
}

// 反转int slice
func ReverseIntSlice(in []int) {
	i, j := 0, len(in)-1
	for i < j {
		in[i], in[j] = in[j], in[i]
		i++
		j--
	}
}

// 反转string slice
func ReverseStringSlice(in []string) {
	i, j := 0, len(in)-1
	for i < j {
		in[i], in[j] = in[j], in[i]
		i++
		j--
	}
}
