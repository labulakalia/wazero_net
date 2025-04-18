package util

// [32 byte][1 err flag][31 length]
func ErrorToUint64(err error) uint64 {
	if err == nil {
		return 0
	}
	errStr := err.Error()
	errLength := uint32(uint32(len(errStr)) | (uint32(1) << uint32(31)))
	errPtr := uint32(StringToPtr(&errStr))
	return Uint32ToUint64(errPtr, uint32(errLength))
}

func Uint64HasError(ret uint64) bool {
	return ret&(uint64((uint32(1) << uint32(31)))) != 0
}

func Uint64ToErrPtrLength(ret uint64) (ptr uint32, length uint32) {
	dataPtr, dataLen := Uint64ToUint32(ret)
	return dataPtr, dataLen ^ (uint32(1) << uint32(31))
}
