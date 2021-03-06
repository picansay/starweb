package srv_found

import (
	"encoding/json"
	"fmt"

	sn "github.com/picansay/starweb/starnode"
)

type SrvFound struct {
	node           *sn.StarNode
	router         *SrvFoundRoute
	remoteNode     map[string]*SrvNodeInfo
	localNode      *SrvNodeInfo
	srvJoinAction  map[string]func(srv *SrvInfo)
	srvLeaveAction map[string]func(srv *SrvInfo)
	group          string
}

func NewSrvFound() *SrvFound {
	node := sn.NewStarNode()
	go node.Server()
	localNode := NewSrvNodeInfo()
	localNode.UUID = node.UUID()
	localNode.Ip = node.Addr()
	localNode.Enable = true

	sf := &SrvFound{
		node:           node,
		localNode:      localNode,
		router:         NewSrvFoundRoute(),
		remoteNode:     make(map[string]*SrvNodeInfo),
		srvJoinAction:  make(map[string]func(srv *SrvInfo)),
		srvLeaveAction: make(map[string]func(srv *SrvInfo)),
		group:          "srv_found",
	}

	sf.listenSrv(sf.group)
	return sf
}

func (self *SrvFound) listenSrv(srv string) error {
	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventExit}, self.srvNodeLeave)

	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventWhisper}, self.srvNodeShout)

	//same to shout
	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventLeave}, self.srvNodeLeave)
	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventShout}, self.srvNodeShout)
	self.node.HandlerFunc(srv, []sn.StarNodeEventType{sn.EventJoin}, self.srvNodeJoin)

	return nil
}

func (self *SrvFound) srvNodeJoin(req *sn.StarNodeRequest) {

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

func (self *SrvFound) srvNodeLeave(req *sn.StarNodeRequest) {
	if node, ok := self.remoteNode[req.Client]; ok {
		for _, srv := range node.Srv {
			self.router.RemoveSrv(srv)
			if actionFun, ok := self.srvLeaveAction[srv.SrvName]; ok {
				actionFun(srv)
			}
		}
		// self.router.RemoveSrv(node)
		delete(self.remoteNode, req.Client)
	}
}

func (self *SrvFound) srvNodeShout(req *sn.StarNodeRequest) {

	node := NewSrvNodeInfo()

	err := json.Unmarshal(req.Msg, node)
	if err != nil {
		fmt.Println("json format:", err)
		return
	}

	fmt.Printf("recv shout msg: %s\n", string(req.Msg))
	fmt.Printf("node: %+v\n", node)

	self.remoteNode[req.Client] = node
	self.handleSrvNode(node)

}

func (self *SrvFound) handleSrvNode(n *SrvNodeInfo) {

	for _, srv := range n.Srv {
		self.router.SetSrv(srv)
		if actionFun, ok := self.srvJoinAction[srv.SrvName]; ok {
			actionFun(srv)
		}
	}

}

func (self *SrvFound) ListenSrv(srvName string, joinFunc func(*SrvInfo), leaveFunc func(*SrvInfo)) {
	if joinFunc != nil {
		self.srvJoinAction[srvName] = joinFunc
	}

	if leaveFunc != nil {
		self.srvLeaveAction[srvName] = joinFunc
	}
}

func (self *SrvFound) Addr() string {
	return self.node.Addr()
}

func (self *SrvFound) Nodes() map[string]*SrvNodeInfo {
	return self.remoteNode
}

func (self *SrvFound) GetSrv(srvName string) ([]*SrvInfo, error) {
	return self.router.GetSrv(srvName)
}

func (self *SrvFound) SetSrv(srvName string, addr string) error {

	self.localNode.SetSrv(srvName, addr)
	self.localNode.SrvEnable = true

	byt, e := json.Marshal(self.localNode)
	if e != nil {
		return e
	}
	// fmt.Printf("localNode: %+v, shout msg:%s\n", self.localNode, string(byt))
	self.node.Shout(self.group, byt)
	return nil
}
