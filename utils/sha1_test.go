package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSha1(t *testing.T) {
	test := []byte("this is my test string")
	data := sha1.Sum(test)
	//fmt.Printf("%#v" ,data)
	hs := hex.EncodeToString(data[:])
	fmt.Println(hs)
}