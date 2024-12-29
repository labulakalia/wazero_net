package util

import (
	"sync"
)

var MemPool = sync.Pool{
	New: func() any {
		return make([]byte, 65535)
	},
}
