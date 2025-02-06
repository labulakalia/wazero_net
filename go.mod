module github.com/labulakalia/wazero_net

go 1.24rc1

require (
	github.com/labulakalia/plugin_api v0.0.12
	github.com/tetratelabs/wazero v1.8.2
)

require google.golang.org/protobuf v1.36.3 // indirect

replace github.com/labulakalia/plugin_api => ../plugin_api/
