package starnode

import (
// "errors"

// "fmt"
)

type StarNodeEnterExitRoute struct {
	r []func(*StarNodeRequest)
}

func NewStarNodeEnterExitRoute() *StarNodeEnterExitRoute {
	return &StarNodeEnterExitRoute{}
}

func (self *StarNodeEnterExitRoute) Handler(req *StarNodeRequest) {
	// if req.Method != EventJoin {
	// 	return
	// }
	// fmt.Println(">>>in join handler")
	// fmt.Println(">>>in join handler", "node:", req.Client, "group:", req.Group)
	// var f func(*StarNode, *StarNodeRequest)
	for _, f := range self.r {
		f(req)
	}
}

func (self *StarNodeEnterExitRoute) HandlerFunc(uri string, method StarNodeEventType, f func(*StarNodeRequest)) error {
	// if method != EventJoin {
	// 	return errors.New("join method error!")
	// }

	self.r = append(self.r, f)
	return nil
}
