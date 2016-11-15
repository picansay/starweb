package starnode

import (
	"errors"
	"fmt"
)

func NotFound(req *StarNodeRequest) {
	fmt.Println("star node api not found! req: ", req)
}

type StarNodeHandler interface {
	Handler(req *StarNodeRequest)
}

type StarNodeRouter interface {
	HandlerFunc(uri string, method StarNodeEventType, f func(*StarNodeRequest)) error
	Handler(req *StarNodeRequest)
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

func (self *StarNodeRoute) Handler(req *StarNodeRequest) {

	if r, ok := self.r[req.Method]; !ok {
		NotFound(req)
	} else {
		r.Handler(req)
	}

}

func (self *StarNodeRoute) HandlerFunc(uri string, method StarNodeEventType, f func(req *StarNodeRequest)) error {

	if r, ok := self.r[method]; ok {

		return r.HandlerFunc(uri, method, f)
	} else {

		return errors.New(fmt.Sprintf("No such method! uri: %s, method: %v, router: %+v", uri, method, r))
	}
}
