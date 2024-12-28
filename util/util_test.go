package util

import (
	"reflect"
	"testing"
)

func TestConvert(t *testing.T) {
	type byteStruct struct {
		bytes []byte
		u32   uint32
	}
	testDatas := []byteStruct{
		{
			bytes: []byte{9, 0, 0, 0},
			u32:   9,
		},
		{
			bytes: []byte{99, 0, 0, 0},
			u32:   99,
		},
		{
			bytes: []byte{231, 3, 0, 0},
			u32:   999,
		},
		{
			bytes: []byte{15, 39, 0, 0},
			u32:   9999,
		},
		{
			bytes: []byte{159, 134, 1, 0},
			u32:   99999,
		},
		{
			bytes: []byte{63, 66, 15, 0},
			u32:   999999,
		},
		{
			bytes: []byte{127, 150, 152, 0},
			u32:   9999999,
		},
		{
			bytes: []byte{255, 224, 245, 5},
			u32:   99999999,
		},
		{
			bytes: []byte{255, 201, 154, 59},
			u32:   999999999,
		},
	}

	for _, v := range testDatas {
		if v.u32 != Byte4ToUint32(v.bytes) {
			t.Fatal("Byte4ToUint32 failed", v.bytes, v.u32)
		}
		if v.u32 != UnsafeByte4ToUint32(v.bytes) {
			t.Fatal("UnsafeByte4ToUint32 failed", v.bytes, v.u32)
		}

		if !reflect.DeepEqual(v.bytes, Uint32toByte4(v.u32)) {
			t.Fatal("Uint32toByte4 failed", v.bytes, v.u32)
		}
		if !reflect.DeepEqual(v.bytes, UnsafeUint32toByte4(v.u32)) {
			t.Fatal("UnsafeUint32toByte4 failed", v.bytes, v.u32)
		}

	}
}

func BenchmarkByte4ToUint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := byte(i % 256)
		Byte4ToUint32([]byte{v, v, v, v})
	}
}

func BenchmarkUnsafeByte4ToUint32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := byte(i % 256)
		UnsafeByte4ToUint32([]byte{v, v, v, v})
	}
}

func BenchmarkUint32ToByte4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uint32toByte4(uint32(i))
	}
}

func BenchmarkUnsafeUint32ToByte4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UnsafeUint32toByte4(uint32(i))
	}
}


func TestUint64ToError(t *testing.T){
	var n1,n2 uint32 = 1021,21121
	n3 := Uint32ToUint64(n1, n2)
	n4,n5 := Uint64ToUint32(n3)
	if n1!=n4 || n2!=n5 {
		t.Fatal("Uint64ToUint32 convert failed")
	}
}
