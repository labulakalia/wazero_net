package http

import (
	"encoding/json"
	"net/http"
	"wazero_net"
	"wazero_net/util"
)


type Transport struct{}

// re impl this interface
func (r *Transport)RoundTrip(req *http.Request) (*http.Response, error) {
	wr := wazero_net.Request{}
	err := wr.ParseHttpRequest(req)
	if err != nil {
		return nil,err
	}
	reqData,err := json.Marshal(wr)
	if err != nil {
		return nil,err
	}
	// TODO wasm limit mem 4096M
	reqDataLen := len(reqData)
	var respLength uint64
	err = util.RetUint64ToError(round_trip(util.BytesToPtr(reqData), uint64(reqDataLen),util.Uint64ToPtr(&respLength)))
	if err != nil {
		return nil,err
	}
	respData := make([]byte,respLength)
	err = util.RetUint64ToError(read_resp(util.BytesToPtr(respData), respLength))
	if err != nil {
		return nil,err
	}

	resp := wazero_net.Response{}
	err = json.Unmarshal(respData, &resp)
	if err != nil {
		return nil,err
	}
	return resp.ToHttpResponse(),nil
}
