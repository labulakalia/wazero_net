package wazero_net

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labulakalia/wazero_net/util"
	"github.com/tetratelabs/wazero/api"
)

func (h *HostNet) client_do(_ context.Context, m api.Module,
	reqPtr, reqLen, respPtr, respLenPtr uint64) uint64 {
	reqBytes, err := util.HostReadBytes(m, uint32(reqPtr), uint32(reqLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	req := &util.Request{}
	err = json.Unmarshal(reqBytes, req)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}

	resp, err := http.DefaultClient.Do(util.ToHttpRequest(req))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}

	rResp := util.Response{
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
		return util.HostErrorToUint64(m, err)
	}
	ptr, err := util.HostWriteBytes(m, respData)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	m.Memory().WriteUint64Le(uint32(respPtr), ptr)
	m.Memory().WriteUint64Le(uint32(respLenPtr), uint64(len(respData)))
	return 0
}
