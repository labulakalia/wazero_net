package wazero_net

import (
	"fmt"

	"github.com/tetratelabs/wazero/api"
)


func ReadBytes(m api.Module, offset, byteCount uint32) ([]byte, error) {
	bytes, ok := m.Memory().Read(offset, byteCount)
	if !ok {
		return nil, fmt.Errorf("read mem failed offset:%d count:%d",offset,byteCount)
	}
	return bytes, nil
}

// func WriteBytes(m api.Module, data []byte) (uint64, error) {
// 	if len(data) == 0 {
// 		return 0, nil
// 	}
// 	ctx := context.Background()
// 	malloc := m.ExportedFunction("malloc")
// 	free := m.ExportedFunction("free")
// 	fmt.Println("malloc",malloc)
// 	result, err := malloc.Call(ctx, uint64(len(data)))
// 	if err != nil {
// 		return 0,fmt.Errorf("malloc mem failed count:%d",len(data))
// 	}
// 	dataPtr := result[0]
// 	defer free.Call(ctx, dataPtr)

// 	ok := m.Memory().Write(uint32(dataPtr), data)
// 	if !ok {
// 		return 0, fmt.Errorf("write mem failed offset:%d count:%d",dataPtr,len(data))
// 	}
// 	return dataPtr, nil
// }


// func ErrorToRetUint64(m api.Module,err error) uint64 {
// 	if err == nil {
// 		return 0
// 	}

// 	errStr := err.Error()
// 	dataPtr,err:= WriteBytes(m, util.StringToBytes(&errStr))
// 	if err != nil {
// 		slog.Error("write bytes failed", "err",err)
// 		panic(err)
// 	}

// 	ret := util.Uint32ToUint64(uint32(dataPtr),uint32(len(errStr)))
// 	var retPtr uint64
// 	ok := m.Memory().WriteUint64Le(uint32(util.Uint64ToPtr(&retPtr)),ret)
// 	if !ok {
// 		slog.Error("write uint64 failed")
// 		panic(err)
// 	}
// 	return retPtr
// }
