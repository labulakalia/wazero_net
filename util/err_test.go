package util

import (
	"errors"
	"testing"
)

func TestErr(t *testing.T) {
	err := errors.New("thisis errr")
	u64 := ErrorToUint64(err)
	if !Uint64HasError(u64) {
		t.Fatal("fail has err")
	}
	ptr, count := Uint64ToErrPtrLength(u64)
	if count != uint32(len(err.Error())) {
		t.Fatal("fail", count, uint32(len(err.Error())))
	}
	if PtrToString(ptr, count) == err.Error() {
		t.Log("success err")
	}
}
