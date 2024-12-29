## WAZERO NET
this project provider net/http for [wazero](https://github.com/tetratelabs/wazero), it not is wazero's official project

>note: unsupport tinygo, because tinygo's gc has some [problem](https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/450#issuecomment-2253729297) for wasm

[go1.24](https://tip.golang.org/doc/go1.24#wasm) will support `go:wasmexport directive` to export function

Dial's Conn can not convert to net.TCPConn,net.UDPConn

## Example

```
cd example
GOOS=wasip1 GOARCH=wasm go build -o http.wasm http.go
GOOS=wasip1 GOARCH=wasm go build -o net.wasm net.go
go run example.go
```

## TODO
- [ ] support listen ip
- [ ] add unit test
