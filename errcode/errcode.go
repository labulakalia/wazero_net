package errcode

import "fmt"

func ToError(code uint64) error {
	if code == 0 {
		return nil
	}
	return fmt.Errorf("errcode: %d", code)
}


const (
	ERR_READ_MEM         uint64 = iota+1
	ERR_WRITE_MEM
	ERR_CONN_DIAL
	ERR_CONN_NOT_EXIST
	ERR_CONN_READ
	ERR_CONN_READ_IO_EOF
	ERR_CONN_WRITE
	ERR_CONN_CLOSE
	ERR_CONN_SET_DEAD_LINE
	ERR_CONN_SET_READ_DEAD_LINE
	ERR_CONN_SET_WRITE_DEAD_LINE
	ERR_LISTENER_NOT_EXIST
	ERR_LISTENER_ACCEPT
	ERR_LISTENER_CLOSE
)
