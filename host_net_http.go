package wazero_net

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/labulakalia/wazero_net/util"
	"github.com/medianexapp/plugin_api/httpclient"
	"github.com/tetratelabs/wazero/api"
)

func (h *HostNet) client_do(ctx context.Context, m api.Module,
	reqPtr, reqLen, respPtr, respLenPtr, timeout uint64) uint64 {
	reqBytes, err := util.HostReadBytes(m, uint32(reqPtr), uint32(reqLen))
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	fmt.Printf("client do %p\n", m)
	slog.Info("client_do", "m", m)
	req := &util.Request{}
	err = json.Unmarshal(reqBytes, req)
	if err != nil {
		return util.HostErrorToUint64(m, err)
	}
	slog.Info("start client do", "url", req.URL, "timeout", timeout)
	resp, err := httpclient.NewBuilder().AutoProxy().SetRequest(util.ToHttpRequest(req)).
		SetTimeout(time.Duration(timeout)).RawResponse()
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
