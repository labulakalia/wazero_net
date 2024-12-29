package http

import (
	"net/http"
)

func init() {
	http.DefaultClient.Transport = &Transport{}
}
