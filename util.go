package utils

import (
	"bufio"
	"crypto/md5"
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
func GetMd5FromString(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

// 计算[]byte的md5值，以32位字符串形式返回
func GetMd5FromBytes(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// 计算文件的md5值，以32位字符串形式返回
func GetMd5FromFile(filename string) (string, error) {
	h := md5.New()
	if err := readFile(filename, h); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 时间戳转换为string显示
func TimestampToString(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// 辅助函数，扫描文件内容并编码到hash.Hash中
func readFile(filename string, h hash.Hash) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(file)
	for s.Scan() {
		h.Write(s.Bytes())
	}

	return s.Err()
}

// 得到一个长度在区间[m, n]内的随机字符串，字母为小写[a, z]
func RandomString(m, n int) string {
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

	return String(bytes)
}

// 不需要拷贝即可返回字符串 s 的 byte slice，但是不能对返回的byte slice做任何修改，否则panic
func Slice(s string) (b []byte) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	pbytes.Cap = pstring.Len
	return b
}

// 不需要拷贝即可返回切片 b 的 string 形式
func String(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstrings := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstrings.Data = pbytes.Data
	pstrings.Len = pbytes.Len
	return s
}
