## WAZERO NET
this project provider net/http for [wazero](https://github.com/tetratelabs/wazero), it not belong to wazero's official project

>note: unsupport tinygo, because tinygo's gc has some [problem](https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/450#issuecomment-2253729297) for wasm

> Dial's Conn can not convert to net.TCPConn,net.UDPConn

## Install
```
go get github.com/labulakalia/wazero_net@v0.0.2
```

## Example
> must use go version >= go1.24, because [go1.24](https://tip.golang.org/doc/go1.24#wasm) will support `go:wasmexport directive` to export function

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

## Todo
- [ ] support listen ip
- [ ] add unit test

## Some Limit or Bug
[sleep](https://github.com/labulakalia/wazero_net/blob/main/example/net.go#L76) [sleep2](https://github.com/labulakalia/wazero_net/blob/main/example/net.go#L34) need call on multi goroutine syscall, or goroutinue can no scheduler,this is bug [example code](https://github.com/golang/go/issues/65178#issuecomment-2565148315)
