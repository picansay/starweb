package starnode

import (
	"errors"
	"fmt"
)

func NotFound(node *StarNode, req *StarNodeRequest) {
	fmt.Println("star node api not found! req: ", req)
}

// func NotFoundHandler() StarNodeHandler { return HandlerFunc(NotFound) }

// type HandlerFunc func(node *StarNode, req *StarNodeRequest) error

// func (f HandlerFunc) Handler(node *StarNode, req *StarNodeRequest) {
// 	f(node, req)
// }

type StarNodeHandler interface {
	Handler(node *StarNode, req *StarNodeRequest)
}

type StarNodeRouter interface {
	HandlerFunc(uri string, method StarNodeEventType, f func(*StarNode, *StarNodeRequest)) error
	Handler(node *StarNode, req *StarNodeRequest)
}

type StarNodeRoute struct {
	r map[StarNodeEventType]StarNodeRouter
}

func NewStarNodeRoute() *StarNodeRoute {
	r := make(map[StarNodeEventType]StarNodeRouter)

	r[EventJoin] = NewStarNodeJoinRoute()
	r[EventLeave] = NewStarNodeLeaveRoute()
	r[EventWhisper] = NewStarNodeWhisperRoute()
	r[EventShout] = NewStarNodeGroupRoute()
	r[EventEnter] = NewStarNodeEnterExitRoute()
	r[EventExit] = NewStarNodeEnterExitRoute()
	return &StarNodeRoute{r: r}
}

func (self *StarNodeRoute) Handler(node *StarNode, req *StarNodeRequest) {

	if r, ok := self.r[req.Method]; !ok {
		NotFound(node, req)
	} else {
		r.Handler(node, req)
	}

}

func (self *StarNodeRoute) HandlerFunc(uri string, method StarNodeEventType, f func(node *StarNode, req *StarNodeRequest)) error {

	if r, ok := self.r[method]; ok {

		return r.HandlerFunc(uri, method, f)
	} else {

		return errors.New(fmt.Sprintf("No such method! uri: %s, method: %v, router: %+v", uri, method, r))
	}
}
