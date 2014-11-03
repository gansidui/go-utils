package utils

import (
	"bytes"
	"testing"
)

func TestUtils1(t *testing.T) {
	// md5
	if BytesToMd5([]byte("hello world")) != StringToMd5("hello world") {
		t.Fatal()
	}

	// PKCS5Padding
	b := PKCS5Padding([]byte("hello"), 8)
	if len(b) != 8 || b[5] != 3 || b[6] != 3 || b[7] != 3 {
		t.Fatal()
	}
	if !bytes.Equal(PKCS5UnPadding(b), []byte("hello")) {
		t.Fatal()
	}

	b = PKCS5Padding([]byte("hello"), 5)
	if len(b) != 10 || b[5] != 5 || b[6] != 5 || b[7] != 5 || b[8] != 5 || b[9] != 5 {
		t.Fatal()
	}
	if !bytes.Equal(PKCS5UnPadding(b), []byte("hello")) {
		t.Fatal()
	}

	// aes
	testAES([]byte("hello你好，这是一个AES加密测试"), []byte("0123456789123456"), t)
	testAES([]byte("hello你好，这是另一个AES加密测试"), []byte("01234567891234560123456789123456"), t)
}

func TestUtils2(t *testing.T) {
	// BigEndian: uint32 <-> []byte, int32 <-> []byte
	if !bytes.Equal(Int32ToBytes(34567), Uint32ToBytes(34567)) {
		t.Fatal()
	}

	b := Int32ToBytes(345345)
	if BytesToInt32(b) != 345345 {
		t.Fatal()
	}

	// *s -> []byte
	s := "hello"
	if !bytes.Equal(StringToByteSlice(&s), []byte(s)) {
		t.Fatal()
	}

	// reverse
	intslice := []int{1, 2, 3}
	ReverseIntSlice(intslice)
	if intslice[0] != 3 || intslice[1] != 2 || intslice[2] != 1 {
		t.Fatal()
	}
	intslice = []int{5, 6, 7, 8}
	ReverseIntSlice(intslice)
	if intslice[0] != 8 || intslice[1] != 7 || intslice[2] != 6 || intslice[3] != 5 {
		t.Fatal()
	}
	stringslice := []string{"hello", "world", "golang"}
	ReverseStringSlice(stringslice)
	if stringslice[0] != "golang" || stringslice[1] != "world" || stringslice[2] != "hello" {
		t.Fatal()
	}

}

func testAES(text, key []byte, t *testing.T) {
	ciphertext, err := AesEncrypt(text, key)
	if err != nil {
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
