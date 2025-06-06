package util

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tetratelabs/wazero/api"
)

func HostReadBytes(m api.Module, offset, byteCount uint32) ([]byte, error) {
	bytes, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		return nil, fmt.Errorf("read mem failed offset %d count %d", offset, byteCount)
	}
	return bytes, nil
}

func HostWriteBytes(m api.Module, data []byte) (uint64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	ctx := context.Background()
	malloc := m.ExportedFunction("malloc")

	result, err := malloc.Call(ctx, uint64(len(data)))
	if err != nil {
		slog.Info("malloc mem failed", "count", len(data), "error", err)
		return 0, fmt.Errorf("malloc mem failed count:%d", len(data))
	}
	dataPtr := result[0]

	ok := m.Memory().Write(uint32(dataPtr), data)
	if !ok {
		return 0, fmt.Errorf("write mem failed offset:%d count:%d", dataPtr, len(data))
	}
	return dataPtr, nil
}

func HostErrorToUint64(m api.Module, err error) uint64 {
	if err == nil {
		return 0
	}
	errStr := err.Error()
	data := StringToBytes(&errStr)
	res, err := HostWriteBytes(m, data)
	if err != nil {
		panic(err)
	}
	return Uint32ToUint64(uint32(res), uint32(len(data)))
}
