package utils

import (
	"bytes"
	"fmt"
	"testing"
)

func TestUtils(t *testing.T) {
	s := "hello world -- 你好，世界"

	// MD5
	if BytesToMd5([]byte(s)) != StringToMd5(s) {
		t.Fatal()
	}

	// *string -> []byte
	if !bytes.Equal(StringToByteSlice(&s), []byte(s)) {
		t.Fatal()
	}

	// AES
	testAES([]byte("hello world"), []byte("0123456789123456"), t)
	testAES([]byte("hello你好，这是一个AES加密测试"), []byte("0123456789123456"), t)
	testAES([]byte("hello你好，这是另一个AES加密测试"), []byte("01234567891234560123456789123456"), t)
	testAES([]byte("hello你好，这又是另一个AES加密测试"), []byte("012345678912345612345678"), t)

	// RSA
	testRSA([]byte("hello world"), publickey, privatekey, t)
	testRSA([]byte("hello你好，这是一个RSA加密测试"), publickey, privatekey, t)
	testRSA([]byte(s), publickey, privatekey, t)
}

func testAES(text, key []byte, t *testing.T) {
	src := make([]byte, len(text))
	copy(src, text)

	ciphertext, err := AesEncrypt(text, key)
	if err != nil {
		t.Fatal()
	}
	if !bytes.Equal(src, text) {
		t.Fatal()
	}

	plaintext, err := AesDecrypt(ciphertext, key)
	if err != nil {
		t.Fatal()
	}
	if !bytes.Equal(text, plaintext) {
		t.Fatal()
	}
}

func testRSA(text, publickey, privatekey []byte, t *testing.T) {
	src := make([]byte, len(text))
	copy(src, text)

	ciphertext, err := RsaEncrypt(text, publickey)
	if err != nil {
		t.Fatal()
	}
	if !bytes.Equal(src, text) {
		t.Fatal()
	}

	src = make([]byte, len(ciphertext))
	copy(src, ciphertext)

	plaintext, err := RsaDecrypt(ciphertext, privatekey)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !bytes.Equal(src, ciphertext) {
		t.Fatal()
	}

	if !bytes.Equal(text, plaintext) {
		t.Fatal()
	}
}

// RSA的公钥和私钥
var privatekey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC+kuVGtO6FT7fmZ1+/q0WvsRloiczIp4toY7OmbozmBkkVgSdl
WPLRO3v5o1cGFkdmJ0wSWJhtrnQ6PH3OZJ9b72HYvh6mlhVHNiMnLkwD3vYli7oK
QUR+DkbU93CIOw/AiF7DhzeEPPAAni6s0HKBo4WT1WUUmciNWulzlfjBnQIDAQAB
AoGAK98HNwrJ6hia/kiH60jTZwm/DqjCYuLnHeXt4n+Koh2KT5AG8LbUV0R5WtO2
YelZEHQ1d/e7R2ykxw9L6uqRIKQRPI5ILUZZYCb9Xrj1b6C0LhFDxTvFGfefTOnl
mCatPIdCmmY9GBR+v+N+IKVTYzRdgY2D6jzcVx2w59fkEIUCQQDqLyw20MYDpaFx
PNjue1VPT7vBDkmKwx6zhqlDxfbAkFLx1acFL9RUS1F9EbsmD3TYweJ17Q9EwmN7
S4w3j/MDAkEA0FOzV0nhV/guuwyXRmnCbbeBEIgtyumYjyTimWhkKy58bUQ+NaXG
Ybbz0A4yPTIWatBaQe+mtDken35RGfYG3wJBAOQ8FVtXHaVwR2eVZdcHXJ1vmA0P
X51djQ5qr4zd4x7Jig0nrR/g/Y8p2MGMBlmRts+KJqvH3pmk2k/P0VhVcwECQQCo
MRtWusgbDL01uMmdSJ93kzK5VSibbRMFZoMn1bchgctlMDaFe4x5sYqQjBWVgI3G
uOZV25UcZg1KOWJi8lXDAkBVhhvvYY4UACWZu6rT1Sw012+SlYgwQviEUvJwRM9h
7xDjIm+Opd9h6JEv2/BpU95hDRVCuqKqqNEJ266tD4H3
-----END RSA PRIVATE KEY-----
`)

var publickey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC+kuVGtO6FT7fmZ1+/q0WvsRlo
iczIp4toY7OmbozmBkkVgSdlWPLRO3v5o1cGFkdmJ0wSWJhtrnQ6PH3OZJ9b72HY
vh6mlhVHNiMnLkwD3vYli7oKQUR+DkbU93CIOw/AiF7DhzeEPPAAni6s0HKBo4WT
1WUUmciNWulzlfjBnQIDAQAB
-----END PUBLIC KEY-----
`)
