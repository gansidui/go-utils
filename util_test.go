package utils

import (
	"bytes"
	"testing"
)

func TestUtil(t *testing.T) {
	s := "hello world -- 你好，世界"

	// md5
	if GetMd5FromString(s) != "5bd20f1bbe2da1d10ca98136c2bd8e26" {
		t.Fatal()
	}

	// slice <--> string
	b := Slice(s)
	if !bytes.Equal(b, []byte(s)) {
		t.Fatal()
	}
	if String(b) != s {
		t.Fatal()
	}
}
