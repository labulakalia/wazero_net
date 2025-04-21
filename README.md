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
tinygo build  -x -target=wasip1 -buildmode=c-shared -o http.wasm http.go
go run example.go http
tinygo build  -x -target=wasip1 -buildmode=c-shared -o net.wasm net.go
go run example.go net
```

## Todo

- [ ] support ip,unix addr
- [ ] add unit test

## Some Limit

https://go.dev/blog/wasi#limitations

Memory Usage is too big,so will switch to tinygo in v2
