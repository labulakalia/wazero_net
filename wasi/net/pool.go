package net

import "sync"


var memPool = sync.Pool{
	New: func() any {
		return make([]byte,65535)
	},
}
