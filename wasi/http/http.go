package http

import (
	"encoding/json"
	"log/slog"

	"github.com/labulakalia/wazero_net/model"
	"github.com/labulakalia/wazero_net/util"
	_ "github.com/labulakalia/wazero_net/wasi/malloc"
)

// re impl this interface
func Do(req *model.Request) (*model.Response, error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		slog.Error("marshal failed", "err", err)
		return nil, err
	}
	// TODO wasm limit mem 4096M
	reqDataLen := len(reqData)

	ret := client_do(util.BytesToPtr(reqData), uint64(reqDataLen))

	dataPtr, dataLen := util.Uint64ToUint32(ret)
	respData := util.PtrToBytes(dataPtr, dataLen)
	slog.Info("re", "respData", string(respData))
	resp := &model.Response{}
	err = json.Unmarshal(respData, resp)
	if err != nil {
		slog.Info("unmarshal failed", "err", err)
		return nil, err
	}
	return resp, nil
}
