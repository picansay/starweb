package starweb

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	sf "github.com/picansay/starweb/srv_found"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

type Server struct {
	ip     string
	ln     *net.TCPListener
	node   *sf.SrvFound
	server *http.Server
	port   int
	r      *mux.Router
}

func NewServer(node *sf.SrvFound) *Server {
	node.ListenSrv("starweb")

	ln, err := net.Listen("tcp", ":0")

	if err != nil {
		return nil
	}

	r := mux.NewRouter()

	server := &http.Server{Handler: r}
	return &Server{
		ip:     node.Addr(),
		r:      r,
		server: server,
		node:   node,
		ln:     ln.(*net.TCPListener),
		port:   ln.(*net.TCPListener).Addr().(*net.TCPAddr).Port,
	}
}

func (self *Server) Addr() string {
	return fmt.Sprintf("http://%s:%d", self.ip, self.port)
}

func (self *Server) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *mux.Route {
	srvName := "starweb"
	if srvName == "" {
		srvName = "starweb"
	}
	self.node.SetSrv(srvName, self.ip, self.port, []string{path})
	return self.r.HandleFunc(path, f)
}

func (self *Server) ListenAndServer() error {
	return self.server.Serve(tcpKeepAliveListener{self.ln})
}
