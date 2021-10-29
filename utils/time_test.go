package utils

import (
	"fmt"
	"testing"
)

func TestTime(t *testing.T)  {
	mt:=StampToTime(1635492336)
	fmt.Println(mt.Unix())
}
