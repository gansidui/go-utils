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
