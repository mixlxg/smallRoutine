package utils

import (
	"crypto/md5"
	"fmt"
)

func MyMd5(password string) string {
	pByte := []byte(password)
	md5Password := md5.Sum(pByte)
	return fmt.Sprintf("%X",md5Password)
}
