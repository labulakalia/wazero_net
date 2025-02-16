package wazero_net

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/labulakalia/wazero_net/model"
	"github.com/labulakalia/wazero_net/util"
	"github.com/tetratelabs/wazero/api"
)

func (h *HostNet) client_do(_ context.Context, m api.Module,
	reqPtr, reqLen uint64) uint64 {
	reqBytes, err := ReadBytes(m, uint32(reqPtr), uint32(reqLen))
	if err != nil {
		return ErrorToUint64(m, err)
	}
	req := &model.Request{}
	err = json.Unmarshal(reqBytes, req)
	if err != nil {
		return ErrorToUint64(m, err)
	}

	resp, err := http.DefaultClient.Do(toHttpRequest(req))
	if err != nil {
		return ErrorToUint64(m, err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ErrorToUint64(m, err)
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
		return ErrorToUint64(m, err)
	}
	malloc := m.ExportedFunction("malloc")
	result, err := malloc.Call(context.Background(), uint64(len(respData)))
	if err != nil {
		return ErrorToUint64(m, err)
	}
	ok := m.Memory().Write(uint32(result[0]), respData)
	if !ok {
		return ErrorToUint64(m, errors.New("write resp data failed"))
	}
	return util.Uint32ToUint64(uint32(result[0]), uint32(len(respData)))
}
