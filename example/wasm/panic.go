//go:build wasip1

package main

import (
	_ "github.com/labulakalia/wazero_net/wasi/http"
)

//go:wasmexport panic_test
func panic_test() {
	panic("panic test")
}

func main() {}
