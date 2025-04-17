## WAZERO NET

this project provider net/http for [wazero](https://github.com/tetratelabs/wazero), it not belong to wazero's official project

> Dial's Conn can not convert to net.TCPConn,net.UDPConn

## Install

```
go get github.com/labulakalia/wazero_net
```

## Example

> must use go version >= go1.24, because [go1.24](https://tip.golang.org/doc/go1.24#wasm) will support `go:wasmexport directive` to export function

> Install Required Go Version
> = 1.24

```
cd example
GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o http.wasm http.go
go run example.go http
GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o net.wasm net.go
go run example.go net
```

## Todo

- [ ] support ip,unix addr
- [ ] add unit test

## Some Limit

https://go.dev/blog/wasi#limitations
