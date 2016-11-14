package starnode

import (
	"errors"
	"fmt"
)

type StarNodeWhisperRoute struct {
	r map[string]func(*StarNode, *StarNodeRequest)
}

func NewStarNodeWhisperRoute() *StarNodeWhisperRoute {
	return &StarNodeWhisperRoute{r: make(map[string]func(*StarNode, *StarNodeRequest))}
}

func (self *StarNodeWhisperRoute) Handler(node *StarNode, req *StarNodeRequest) {
	if req.Method != EventWhisper {
		return
	}
	// fmt.Println(">>>in Whisper handler")
	fmt.Println(self.r)
	fmt.Println(">>>in Whisper handler", "node:", req.Client, "uri:", req.Uri)
	if req.Uri == "" {
		return
	}

	var f func(*StarNode, *StarNodeRequest)

	if hf, ok := self.r[req.Uri]; !ok {
		f = NotFound
	} else {
		f = hf
	}

	f(node, req)
}

func (self *StarNodeWhisperRoute) HandlerFunc(uri string, method StarNodeEventType, f func(*StarNode, *StarNodeRequest)) error {
	if method != EventWhisper {
		return errors.New("EventWhisper method error!")
	}

	if uri == "" {
		return errors.New("EventWhisper uri error!")
	}
	if _, ok := self.r[uri]; ok {
		return nil
	} else {
		self.r[uri] = f
	}
	return nil
}
