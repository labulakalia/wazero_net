## WAZERO NET
this project provider net/http for [wazero](https://github.com/tetratelabs/wazero), it not is wazero's official project

>note: unsupport tinygo, because tinygo's gc has some [problem](https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/450#issuecomment-2253729297) for wasm

[go1.24](https://tip.golang.org/doc/go1.24#wasm) will support `go:wasmexport directive` to export function

Dial's Conn can not convert to net.TCPConn,net.UDPConn

## Example
> must use version >= go1.24

> Install Required Go Version
```
go install golang.org/dl/go1.24rc1@latest
go1.24rc1 download
```
```
cd example
GOOS=wasip1 GOARCH=wasm go1.24rc1 build -buildmode=c-shared -o http.wasm http.go
GOOS=wasip1 GOARCH=wasm go1.24rc1 build -buildmode=c-shared -o net.wasm net.go
go1.24rc1 run example.go
```

## TODO
- [ ] support listen ip
- [ ] add unit test
