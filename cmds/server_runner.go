package cmds

import (
	"net/http"
)

func RunServer(handler http.Handler, addr string) error {
	return http.ListenAndServe(addr, handler)
}
