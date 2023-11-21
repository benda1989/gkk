package tool

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomStr(n int, allowedChars ...[]rune) string {
	var letter = letters
	if len(allowedChars) != 0 {
		letter = allowedChars[0]
	}
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func RandomStrLetter(n int) string {
	return RandomStr(n, []rune("abcdefghijklmnopqrstuvwxyz"))
}

func RandomCode(ns ...int) string {
	n := 6
	if len(ns) > 0 {
		n = ns[0]
	}
	return RandomStr(n, []rune("0123456789"))
}

func SHA1(str any) string {
	o := sha1.New()
	switch str.(type) {
	case []byte:
		o.Write(str.([]byte))
	case string:
		o.Write([]byte(str.(string)))
	default:
		fmt.Println("SHA1() arg error")
	}
	return hex.EncodeToString(o.Sum(nil))
}

func Base64(str []byte) (re []byte) {
	e64 := base64.StdEncoding
	re = make([]byte, e64.EncodedLen(len(str)))
	e64.Encode(re, str)
	return
}

func MD5(str any) string {
	_md5 := md5.New()
	switch str.(type) {
	case []byte:
		_md5.Write(str.([]byte))
	case string:
		_md5.Write([]byte(str.(string)))
	default:
		fmt.Println("MD5() arg error")
	}
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func MD5Struct(v any) string {
	payload, _ := msgpack.Marshal(v)
	return MD5(payload)
}
