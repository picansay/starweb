package starweb

import (
	// "errors"
	// "fmt"

	// "encoding/json"
	// "log"
	// "net"
	"net/http"
	// "time"

	"github.com/gorilla/mux"
	// "github.com/zeromq/gyre"
)

type StarWebRouter struct {
	r *mux.Router
}

func NewStarWebRouter() *StarWebRouter {
	return &StarWebRouter{mux.NewRouter()}
}
func (self StarWebRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	self.r.ServeHTTP(w, req)
}

func (self StarWebRouter) SetNotFoundHandler(handler http.Handler) {
	self.r.NotFoundHandler = handler
}

func (self *StarWebRouter) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *mux.Route {
	return self.r.HandleFunc(path, f)
}
