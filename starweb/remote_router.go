package starweb

import (
	// "errors"
	"fmt"

	"encoding/json"
	// "log"
	// "net"
	"net/http"
	// "time"

	"github.com/gorilla/mux"
	// "github.com/zeromq/gyre"

	sn "github.com/picansay/starweb/starnode"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Remote 404 page not found! ", http.StatusNotFound)
}

func NotFoundHandler() http.Handler { return http.HandlerFunc(NotFound) }

type RemoteHandler interface {
}

type RemoteHandle struct {
	// http.Handler
	Uri   string                     `json:"uri"`
	Nodes map[string]*RemoteNodeInfo `json:"nodes"`
}

func NewRemoteHandle() *RemoteHandle {
	return &RemoteHandle{
		Nodes: make(map[string]*RemoteNodeInfo),
	}
}

func (self *RemoteHandle) AddNode(rni *RemoteNodeInfo) {
	fmt.Println("1>>>>>> in add node: ", "rh:", self, ", nodes: ", self.Nodes)
	// self.Nodes = append(self.Nodes, node)
	self.Nodes[rni.UUID] = rni

	fmt.Println("2>>>>>> in add node: ", "rh:", self, ", nodes: ", self.Nodes)
}

func (self *RemoteHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	b, _ := json.Marshal(self)
	w.Write(b)
}

type RemoteNodeInfo struct {
	UUID   string   `json:"uuid"`
	Addr   string   `json:"addr"`
	Uris   []string `json:"uris"`
	Enable bool     `json:"enable"`
}

type StarWebRemoteRouter struct {
	r            *mux.Router
	node         *sn.StarNode
	remote_nodes map[string]*RemoteNodeInfo
	ps           map[string]RemoteHandler
	node_ctl     chan string
}

func NewStarWebRemoteRouter(node *sn.StarNode) *StarWebRemoteRouter {

	// node := sn.NewStarNode()
	remote_nodes := make(map[string]*RemoteNodeInfo)
	r := mux.NewRouter()
	r.NotFoundHandler = NotFoundHandler()
	return &StarWebRemoteRouter{
		r:            r,
		node:         node,
		remote_nodes: remote_nodes,
		ps:           make(map[string]RemoteHandler),
		node_ctl:     make(chan string),
	}
}
func (self *StarWebRemoteRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	self.r.ServeHTTP(w, req)
}

func (self *StarWebRemoteRouter) HandlerRemoteNodeInfo(msg []byte) {
	var rni RemoteNodeInfo
	err := json.Unmarshal(msg, &rni)
	if err != nil {
		fmt.Println("json format:", err)
	}

	self.remote_nodes[rni.UUID] = &rni
	fmt.Println("recv:", rni)
	for _, uri := range rni.Uris {
		var handler *RemoteHandle

		if h, ok := self.ps[uri]; ok {
			handler = h.(*RemoteHandle)
			handler.AddNode(&rni)
		} else {
			handler = NewRemoteHandle()
			handler.Uri = uri
			handler.AddNode(&rni)

			self.ps[uri] = handler
		}
		fmt.Printf("handler : %p\n", handler)
		self.r.Handle(uri, handler).Methods("GET")
	}
}

// func (self *StarWebRemoteRouter) AutoWebNodeFound(node *sn.StarNode, req *sn.StarNodeRequest) {

// }

// func (self *StarWebRemoteRouter) HandleFunc(path string, f func(http.ResponseWriter,
// 	*http.Request)) *mux.Route {
// 	fmt.Println("remote handler: ", path, ",f:", f)
// 	return self.r.HandleFunc(path, f)
// }
