package starweb

import (
	// "encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	sf "github.com/picansay/starweb/srv_found"
	// sn "github.com/picansay/starweb/starnode"
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

type LocalUris struct {
	Uris []string `json:"uris"`
	Addr string   `json:"addr"`
}

type StarWebNodeInfo struct {
	UUID string `json:"uuid"`
	Addr string `json:"addr"`
}

type StarWeb struct {
	node *sf.SrvFound
	// ip             string
	port   int
	r      *StarWebRouter
	server *http.Server
	ln     *net.TCPListener
	// remote_handler *StarWebRemoteRouter
	// ps             []string
	// work           *StarWebWorker
}

func NewStarWeb(node *sf.SrvFound) *StarWeb {
	addr := ":0"

	var ps []string
	ln, _ := net.Listen("tcp", addr)

	remote_handler := NewStarWebRemoteRouter(node)
	// remote_handler.SyncRouter()
	r := NewStarWebRouter()
	r.SetNotFoundHandler(remote_handler)

	server := &http.Server{Handler: r}
	web := &StarWeb{
		node:   node,
		r:      r,
		ln:     ln.(*net.TCPListener),
		server: server,
		// remote_handler: remote_handler,
		// ps: ps,
		// ip:             node.Addr(),
		port: ln.(*net.TCPListener).Addr().(*net.TCPAddr).Port,
	}

	return web
}

func (self *StarWeb) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *mux.Route {
	// ps := *self.ps.([]string)
	// self.ps = append(self.ps, path)

	return self.r.HandleFunc(path, f)
}

// func (self *StarWeb) Addr() string {
// 	return fmt.Sprintf("http://%s:%d", self.ip, self.port)
// }

// func (self *StarWeb) UUID() string {
// 	return self.node.UUID()
// }

func (self *StarWeb) onLine() error {

	work := NewStarWebWorker(self, self.node)
	self.work = work
	return nil
}

func (self *StarWeb) Loop() error {
	fmt.Printf("=======local path===========\n")
	for _, path := range self.ps {
		fmt.Printf("path: [%s]\n", path)
	}
	// self.node.Whisper()

	// self.remote_handler.SetLocalUrisToRemote(self.ps)

	self.onLine()
	return self.server.Serve(tcpKeepAliveListener{self.ln})
}
