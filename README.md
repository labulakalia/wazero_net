## WAZERO NET

this project provider net/http for [wazero](https://github.com/tetratelabs/wazero), it not belong to wazero's official project

## Install Latest TinyGo

[Install](https://tinygo.org/getting-started/install/)

## Import Package

```
go get github.com/labulakalia/wazero_net
```

## Example

```
cd example
tinygo build  -x -target=wasip1 -buildmode=c-shared -o http.wasm wasm/http.go
go run example.go http
tinygo build  -x -target=wasip1 -buildmode=c-shared -o net.wasm wasm/net.go
go run example.go net
```

## Todo

- [ ] support ip,unix addr
- [ ] add unit test
