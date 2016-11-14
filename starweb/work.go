package starweb

// import (
// 	"encoding/json"
// 	"fmt"

// 	sn "github.com/picansay/starweb/starnode"
// )

// type StarWebWorker struct {
// 	web  *StarWeb
// 	node *sn.StarNode
// }

// func NewStarWebWorker(web *StarWeb, node *sn.StarNode) *StarWebWorker {
// 	self := &StarWebWorker{
// 		web:  web,
// 		node: node,
// 	}

// 	self.node.HandlerFunc("/starweb/node", []sn.StarNodeEventType{sn.EventJoin}, self.WebNode)

// 	self.node.HandlerFunc("/starweb/node/apis", []sn.StarNodeEventType{sn.EventWhisper}, self.WebNodeApis)
// 	return self
// }

// func (self *StarWebWorker) WebNode(node *sn.StarNode, req *sn.StarNodeRequest) {

// 	var rni RemoteNodeInfo

// 	rni.Addr = self.web.Addr()
// 	rni.UUID = self.web.UUID()
// 	rni.Uris = self.web.ps

// 	fmt.Println(">>>>recv join and whisper to node-->", ",req:", req, ", whisper msg: ", rni)
// 	byt, _ := json.Marshal(&rni)
// 	self.node.Whisper(req.Client, "/starweb/node/apis", byt)
// }

// func (self *StarWebWorker) WebNodeApis(node *sn.StarNode, req *sn.StarNodeRequest) {

// 	self.web.remote_handler.HandlerRemoteNodeInfo(req.Msg)

// }
