package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/labulakalia/wazero_net/util"
	_ "github.com/labulakalia/wazero_net/wasi/malloc"
)

// re impl this interface
func Do(req *http.Request) (*http.Response, error) {
	modelReq, err := util.ToModelRequest(req)
	if err != nil {
		slog.Error("to model request failed", "err", err)
		return nil, err
	}

	reqData, err := json.Marshal(modelReq)
	if err != nil {
		slog.Error("marshal failed", "err", err)
		return nil, err
	}

	// TODO wasm limit mem 4096M
	reqDataLen := len(reqData)

	ret := client_do(util.BytesToPtr(reqData), uint64(reqDataLen))

	dataPtr, dataLen := util.Uint64ToUint32(ret)
	respData := util.PtrToBytes(dataPtr, dataLen)
	resp := &util.Response{}
	err = json.Unmarshal(respData, resp)
	if err != nil {
		slog.Error("unmarshal failed", "err", err)
		return nil, err
	}
	return util.ToHttpResponse(resp), nil
}
