package wazero_net

import (
	"context"
	"fmt"

	"github.com/labulakalia/wazero_net/util"
	"github.com/tetratelabs/wazero/api"
)

func ReadBytes(m api.Module, offset, byteCount uint32) ([]byte, error) {
	bytes, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		return nil, fmt.Errorf("read mem failed offset %d count %d", offset, byteCount)
	}
	return bytes, nil
}

func WriteBytes(m api.Module, data []byte) (uint64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	ctx := context.Background()
	malloc := m.ExportedFunction("malloc")

	result, err := malloc.Call(ctx, uint64(len(data)))
	if err != nil {
		return 0, fmt.Errorf("malloc mem failed count:%d", len(data))
	}
	dataPtr := result[0]

	ok := m.Memory().Write(uint32(dataPtr), data)
	if !ok {
		return 0, fmt.Errorf("write mem failed offset:%d count:%d", dataPtr, len(data))
	}
	return dataPtr, nil
}

func ErrorToUint64(m api.Module, err error) uint64 {
	if err == nil {
		return 0
	}
	errStr := err.Error()
	data := util.StringToBytes(&errStr)
	res, err := WriteBytes(m, data)
	if err != nil {
		panic(err)
	}
	return util.Uint32ToUint64(uint32(res), uint32(len(data)))
}
