package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOs(t *testing.T)  {
	var basePath string
	currentPath, err := os.Getwd()
	if err !=nil{
		t.Fatalf("获取当前文件路径失败，错误信息：%v", err.Error())
	}
	if strings.Contains(currentPath,"utils"){

		basePath = filepath.Dir(currentPath)
	}
	fmt.Println(currentPath)
	fmt.Println(basePath)
}
