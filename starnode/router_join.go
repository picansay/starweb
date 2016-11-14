package starnode

import (
	"errors"

	"fmt"
)

type StarNodeJoinRoute struct {
	r []func(*StarNode, *StarNodeRequest)
}

func NewStarNodeJoinRoute() *StarNodeJoinRoute {
	return &StarNodeJoinRoute{}
}

func (self *StarNodeJoinRoute) Handler(node *StarNode, req *StarNodeRequest) {
	if req.Method != EventJoin {
		return
	}
	// fmt.Println(">>>in join handler")
	fmt.Println(">>>in join handler", "node:", req.Client, "group:", req.Group)
	// var f func(*StarNode, *StarNodeRequest)
	for _, f := range self.r {
		f(node, req)
	}
}

func (self *StarNodeJoinRoute) HandlerFunc(uri string, method StarNodeEventType, f func(*StarNode, *StarNodeRequest)) error {
	if method != EventJoin {
		return errors.New("join method error!")
	}

	self.r = append(self.r, f)
	return nil
}
