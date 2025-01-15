package util

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMempoll(t *testing.T) {
	data := MemPool.Get().([]byte)
	data = bytes.Repeat([]byte{0}, len(data))
	fmt.Println(data)
}
