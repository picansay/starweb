package srv_found

import (
	"encoding/json"
	"fmt"

	sn "github.com/picansay/starweb/starnode"
)

type SrvFound struct {
	node       *sn.StarNode
	router     *SrvFoundRoute
	remoteNode map[string]*SrvNodeInfo
	localNode  *SrvNodeInfo
}

func NewSrvFound() *SrvFound {
	node := sn.NewStarNode()
	go node.Server()
	localNode := NewSrvNodeInfo()
	localNode.UUID = node.UUID()
	localNode.Ip = node.Addr()
	localNode.Enable = true

	return &SrvFound{
		node:       node,
		localNode:  localNode,
		router:     NewSrvFoundRoute(),
		remoteNode: make(map[string]*SrvNodeInfo),
	}
}

func (self *SrvFound) ListenSrv(srv string) error {
	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventExit}, self.srvNodeLeave)

	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventWhisper}, self.srvNodeShout)

	//same to shout
	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventLeave}, self.srvNodeLeave)
	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventShout}, self.srvNodeShout)
	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventJoin}, self.srvNodeJoin)

	return nil
}

func (self *SrvFound) srvNodeJoin(n *sn.StarNode, req *sn.StarNodeRequest) {

	if node, ok := self.remoteNode[req.Client]; ok {
		node.Enable = true
	} else {
		node = NewSrvNodeInfo()
		node.UUID = req.Client
		node.Ip = req.Addr
		node.Enable = true

		self.remoteNode[req.Client] = node
	}

	if self.localNode.SrvEnable {
		byt, _ := json.Marshal(self.localNode)
		self.node.Whisper(req.Client, req.Group, byt)
	}
}

func (self *SrvFound) srvNodeLeave(n *sn.StarNode, req *sn.StarNodeRequest) {
	if node, ok := self.remoteNode[req.Client]; ok {
		// node.Enable = false
		// self.remoteNode[req.Client]

		self.router.RemoveSrvNode(node)
		delete(self.remoteNode, req.Client)
	}
}

func (self *SrvFound) srvNodeShout(n *sn.StarNode, req *sn.StarNodeRequest) {

	node := NewSrvNodeInfo()

	err := json.Unmarshal(req.Msg, node)
	if err != nil {
		fmt.Println("json format:", err)
		return
	}

	fmt.Printf("recv shout msg: %s\n", string(req.Msg))
	fmt.Printf("node: %+v\n", node)
	self.router.HandleSrvNode(node)
	self.remoteNode[req.Client] = node
}

func (self *SrvFound) GetSrv(srvName string, uri string) ([]*SrvInfo, error) {
	return self.router.SrvNode(srvName, uri)
}

func (self *SrvFound) Addr() string {
	return self.node.Addr()
}

func (self *SrvFound) Nodes() map[string]*SrvNodeInfo {
	return self.remoteNode
}

//
func (self *SrvFound) SetSrv(srvName string, ip string, port int, uri []string) error {

	self.localNode.SetSrv(srvName, ip, port, uri)
	self.localNode.SrvEnable = true

	byt, e := json.Marshal(self.localNode)
	if e != nil {
		return e
	}
	fmt.Printf("localNode: %+v, shout msg:%s\n", self.localNode, string(byt))
	self.node.Shout(srvName, byt)
	return nil
}
