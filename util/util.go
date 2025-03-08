package util

import (
	"errors"
	"unsafe"
)

func BytesToUint32Arr(bytes []byte) []uint32 {
	res := []uint32{}
	for i := 0; i < len(bytes); i += 4 {
		res = append(res, Byte4ToUint32(bytes[i:i+4]))
	}
	return res
}

func Uint32toByte4(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

func Byte4ToUint32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}

func UnsafeUint32toByte4(val uint32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&val)), 4)
}

func UnsafeByte4ToUint32(val []byte) uint32 {
	return *(*uint32)(unsafe.Pointer(&val[0]))
}

// https://colobu.com/2022/09/06/string-byte-convertion/
func BytesToString(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

func StringToBytes(s *string) []byte {
	return unsafe.Slice(unsafe.StringData(*s), len(*s))
}

func StringToPtr(s *string) uint64 {
	return uint64(uintptr(unsafe.Pointer(unsafe.StringData(*s))))
}

func PtrToString(ptr uint32, ptrLen uint32) string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(ptr))), int(ptrLen))
}

func BytesToPtr(bytes []byte) uint64 {
	if len(bytes) == 0 {
		return 0
	}
	return uint64(uintptr(unsafe.Pointer(&bytes[0])))
}

func PtrToBytes(ptr uint32, ptrLen uint32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), int(ptrLen))
}

func Uint64ToPtr(u64 *uint64) uint64 {
	return uint64(uintptr(unsafe.Pointer(u64)))
}

func Uint32ToPtr(u32 *uint32) uint64 {
	return uint64(uintptr(unsafe.Pointer(u32)))
}

func Uint64ToUint32(u64 uint64) (uint32, uint32) {
	return uint32(u64 >> 32), uint32(u64 & (1<<32 - 1))
}

func Uint32ToUint64(n1u32 uint32, n2u32 uint32) uint64 {
	return uint64(uint64(n1u32)<<32) + uint64(n2u32)
}

func RetUint64ToError(u64 uint64) error {
	if u64 == 0 {
		return nil
	}
	return errors.New(PtrToString(Uint64ToUint32(u64)))
}
