//go:build wasip1

package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/labulakalia/wazero_net/util"
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
	var (
		respPtr uint64
		respLen uint64
	)
	ret := _client_do(util.BytesToPtr(reqData), uint64(reqDataLen), util.Uint64ToPtr(&respPtr), util.Uint64ToPtr(&respLen))
	if ret != 0 {
		return nil, util.RetUint64ToError(ret)
	}

	respData := util.PtrToBytes(uint32(respPtr), uint32(respLen))
	resp := &util.Response{}
	err = json.Unmarshal(respData, resp)
	if err != nil {
		slog.Error("unmarshal failed", "err", err)
		return nil, err
	}
	return util.ToHttpResponse(resp), nil
}
