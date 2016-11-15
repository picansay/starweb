package starnode

import (
	"errors"
	// "fmt"
)

type StarNodeGroupRoute struct {
	r map[string]func(req *StarNodeRequest)
}

func NewStarNodeGroupRoute() *StarNodeGroupRoute {
	return &StarNodeGroupRoute{r: make(map[string]func(*StarNodeRequest))}
}

func (self *StarNodeGroupRoute) Handler(req *StarNodeRequest) {
	// fmt.Println("start handler--")
	if req.Method != EventShout {
		return
	}
	// fmt.Println(">>>in Group handler", "node:", req.Client, "group:", req.Group)
	if req.Group == "" {
		return
	}

	var f func(req *StarNodeRequest)

	if hf, ok := self.r[req.Group]; ok {
		f = hf
	} else {
		f = NotFound
	}

	f(req)
}

func (self *StarNodeGroupRoute) HandlerFunc(group string, method StarNodeEventType, f func(req *StarNodeRequest)) error {
	if method != EventShout {
		return errors.New("shout method error!")
	}

	if group == "" {
		return errors.New("EventShout uri error!")
	}
	if _, ok := self.r[group]; ok {
		return nil
	} else {
		self.r[group] = f
	}
	return nil
}
