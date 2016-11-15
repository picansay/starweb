package starnode

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/zeromq/gyre"
)

type StarNodeEventType interface{}

const (
	EventEnter   = gyre.EventEnter
	EventJoin    = gyre.EventJoin
	EventLeave   = gyre.EventLeave
	EventExit    = gyre.EventExit
	EventWhisper = gyre.EventWhisper
	EventShout   = gyre.EventShout
)

type StarNode struct {
	node *gyre.Gyre
	r    StarNodeRouter
	// ready       bool
	shoutInput  chan ShoutSendMsg
	wisperInput chan WhisperSendMsg
}

type WhisperMsg struct {
	Uri string `json:"uri"`
	Msg []byte `json:"msg"`
}

func NewStarNode() *StarNode {
	node, err := gyre.New()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	// addr, _ := node.Addr()
	err = node.Start()
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	fmt.Printf("sn node: %s\n", node.UUID())
	shoutInput := make(chan ShoutSendMsg, 1)
	wisperInput := make(chan WhisperSendMsg, 1)

	return &StarNode{
		node:        node,
		r:           NewStarNodeRoute(),
		shoutInput:  shoutInput,
		wisperInput: wisperInput,
	}
}

func (self *StarNode) Addr() string {
	addr, _ := self.node.Addr()
	return addr
}

type ShoutSendMsg struct {
	group string
	msg   []byte
}

func (self *StarNode) Shout(group string, msg []byte) error {
	// if !self.ready {
	// 	for {
	// 		if self.ready {
	// 			fmt.Println("starnode ready!")
	// 			break
	// 		}

	// 		time.Sleep(time.Second * 1)
	// 	}
	// }
	var m ShoutSendMsg
	m.group = group
	m.msg = msg
	self.shoutInput <- m
	// fmt.Println("star node shout: ", group, string(msg))
	return nil
	// return self.node.Shout(group, msg)
}

type WhisperSendMsg struct {
	peer string
	uri  string
	msg  []byte
}

func (self *StarNode) Whisper(peer string, uri string, msg []byte) error {
	if peer == self.node.UUID() {
		return errors.New("Can't whisper self!")
	}
	var m WhisperSendMsg

	m.peer = peer
	m.uri = uri
	// m.msg = msg

	var wmsg WhisperMsg

	wmsg.Uri = uri
	wmsg.Msg = msg

	if byt, e := json.Marshal(wmsg); e != nil {
		return e
	} else {
		m.msg = byt
	}

	fmt.Println("whisper--1->", m)
	self.wisperInput <- m
	fmt.Println("whisper--1-end>", m)
	return nil
}

func (self *StarNode) UUID() string {
	return self.node.UUID()
}

func (self *StarNode) HandlerFunc(uri string, methods []StarNodeEventType, f func(req *StarNodeRequest)) error {
	var e error
	fmt.Println("hf: ", uri, methods, f)
	for _, method := range methods {
		fmt.Printf("type: %T\n", method)
		e = self.r.HandlerFunc(uri, method, f)
		if e != nil {
			fmt.Println("error:", e)
			return e
		}
		fmt.Println("handler method:", method, ", uri:", uri)
		// fmt.Println(method, method.(Type))
		if method.(gyre.EventType) == EventShout || method.(gyre.EventType) == EventJoin {
			fmt.Println("====>>>>", self.node.UUID, ",im join to: ", uri)
			time.Sleep(time.Second * 3)
			self.node.Join(uri)
		}
	}
	return nil
}
func (self *StarNode) eventRecv(e *gyre.Event) {
	var req StarNodeRequest
	req.Client = e.Sender()
	req.Addr = e.Addr()
	req.Msg = e.Msg()
	req.Method = e.Type()
	fmt.Println("---->recv method: ", e.Type(), ", uuid: ", e.Sender(), "req:", req)

	switch e.Type() {

	case gyre.EventEnter:

	case gyre.EventLeave:
		req.Group = e.Group()
	case gyre.EventJoin:
		fmt.Println("recv method: ", e.Type())
		req.Group = e.Group()
		// self.node.Shout(e.Group(), []byte("aaaaaaaaaaaaaaa"))
	case gyre.EventShout:
		req.Group = e.Group()
		// fmt.Println("recv shout!")
		// req.Uri = e.Group()
		// fmt.Printf("recv group: %+v\n", req)
	case gyre.EventWhisper:
		fmt.Println("recv method: ", e.Type())
		var wmsg WhisperMsg
		if err := json.Unmarshal(req.Msg, &wmsg); err != nil {
			//TODO:error
		} else {
			req.Uri = wmsg.Uri
			req.Msg = wmsg.Msg
		}
		// fmt.Printf("recv whisper: %+v\n", req)

	}
	// fmt.Printf("recv total: %+v\n", req)
	// node := self
	self.r.Handler(&req)
}
func (self *StarNode) Server() {
	fmt.Println("set starnode ready!")
	// self.ready = true
	// time.Sleep(10)

	for {
		select {
		case e := <-self.node.Events():
			self.eventRecv(e)
		case wmsg := <-self.wisperInput:
			fmt.Println("starnode whisper eeeee: ", wmsg)
			self.node.Whisper(wmsg.peer, wmsg.msg)
		case smsg := <-self.shoutInput:
			fmt.Println("starnode shout eeeee: ", smsg)
			self.node.Shout(smsg.group, smsg.msg)
		}
	}
}
