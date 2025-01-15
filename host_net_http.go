package wazero_net

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/labulakalia/wazero_net/errcode"
	"github.com/labulakalia/wazero_net/model"
	"github.com/labulakalia/wazero_net/util"
	"github.com/tetratelabs/wazero/api"
)

func (h *HostNet) client_do(_ context.Context, m api.Module,
	reqPtr, reqLen uint64) uint64 {
	reqBytes, err := ReadBytes(m, uint32(reqPtr), uint32(reqLen))
	if err != nil {
		slog.Error("listener not found", "err", err)
		return errcode.ERR_READ_MEM
	}
	req := &model.Request{}
	err = json.Unmarshal(reqBytes, req)
	if err != nil {
		slog.Error("json unmarshal failed", "err", err)
		return errcode.ERR_READ_MEM
	}

	resp, err := http.DefaultClient.Do(toHttpRequest(req))
	if err != nil {
		slog.Error("client roundtrip failed", "err", err)
		return errcode.ERR_READ_MEM
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("read all failed", "err", err)
		return errcode.ERR_READ_MEM
	}

	rResp := model.Response{
		StatusCode:    resp.StatusCode,
		Proto:         resp.Proto,
		ProtoMajor:    resp.ProtoMajor,
		ProtoMinor:    resp.ProtoMinor,
		Header:        resp.Header,
		Body:          respBytes,
		ContentLength: resp.ContentLength,
	}
	respData, err := json.Marshal(rResp)
	if err != nil {
		slog.Error("json marshal failed", "err", err)
		return errcode.ERR_READ_MEM
	}
	malloc := m.ExportedFunction("malloc")
	result, err := malloc.Call(context.Background(), uint64(len(respData)))
	if err != nil {
		slog.Error("malloc call failed", "err", err)
		return errcode.ERR_WRITE_MEM
	}
	m.Memory().Write(uint32(result[0]), respData)
	return util.Uint32ToUint64(uint32(result[0]), uint32(len(respData)))
}
