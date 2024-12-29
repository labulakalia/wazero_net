## WAZERO NET
this project provider net/http for [wazero](https://github.com/tetratelabs/wazero), it not is wazero's official project

>note: unsupport tinygo, because tinygo's gc has some [problem](https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/450#issuecomment-2253729297) for wasm

[go1.24](https://tip.golang.org/doc/go1.24#wasm) will support `go:wasmexport directive` to export function

## Usage
go to [example](https://github.com/labulakalia/wazero_net/tree/main/example)

## Next Plan
- [ ] add unit test for net/http
- [ ] add ci for all platform,like windows,macos,linux,andriod,ios
