package srv_found

import (
	"errors"
	"fmt"
)

type SrvGroup struct {
	srv map[string]*SrvInfo
}

func NewSrvGroup() *SrvGroup {
	return &SrvGroup{srv: make(map[string]*SrvInfo)}
}

func (self *SrvGroup) SetSrv(srv *SrvInfo) error {
	// self.srv = append(self.srv, srv)
	self.srv[srv.UUID] = srv
	return nil
}

func (self *SrvGroup) RemoveSrv(srv *SrvInfo) error {
	// self.srv = append(self.srv, srv)
	delete(self.srv, srv.UUID)
	return nil
}

func (self *SrvGroup) GetSrv(srv string) ([]*SrvInfo, error) {
	var srvArry []*SrvInfo

	for _, v := range self.srv {
		srvArry = append(srvArry, v)
	}
	return srvArry, nil
}

type SrvFoundRoute struct {
	router map[string]*SrvGroup
	// nodes []*SrvNodeInfo
}

func NewSrvFoundRoute() *SrvFoundRoute {
	return &SrvFoundRoute{
		router: make(map[string]*SrvGroup),
	}
}

func (self *SrvFoundRoute) SetSrv(srv *SrvInfo) error {
	var group *SrvGroup
	if g, ok := self.router[srv.SrvName]; ok {
		group = g
	} else {
		group = NewSrvGroup()
		self.router[srv.SrvName] = group
	}

	group.SetSrv(srv)
	return nil
}

func (self *SrvFoundRoute) RemoveSrv(srv *SrvInfo) error {
	if g, ok := self.router[srv.SrvName]; ok {
		return g.RemoveSrv(srv)
	}
	return nil
}

func (self *SrvFoundRoute) GetSrv(srvName string) ([]*SrvInfo, error) {
	if g, ok := self.router[srvName]; ok {
		return g.GetSrv(srvName)
	}

	return nil, errors.New(fmt.Sprintf("Not Found Such Srv! [%s]", srvName))
}
