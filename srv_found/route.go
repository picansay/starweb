package srv_found

import (
	"errors"
	"fmt"
)

type SigUriRoute struct {
	//k: node_uuid , v: srvInfo
	srv map[string]*SrvInfo
}

func NewSigUriRoute() *SigUriRoute {
	return &SigUriRoute{srv: make(map[string]*SrvInfo)}
}

func (self *SigUriRoute) HandleSrvNode(srv *SrvInfo) error {
	// self.srv = append(self.srv, srv)
	self.srv[srv.UUID] = srv
	return nil
}

func (self *SigUriRoute) SrvNode(srv string, uri string) ([]*SrvInfo, error) {
	var srvArry []*SrvInfo

	for _, v := range self.srv {
		srvArry = append(srvArry, v)
	}
	return srvArry, nil
}

func (self *SigUriRoute) RemoveSrvNode(srv *SrvInfo) error {
	delete(self.srv, srv.UUID)
	return nil
}

type UriRoute struct {
	router map[string]*SigUriRoute
}

func NewUriRoute() *UriRoute {
	return &UriRoute{
		router: make(map[string]*SigUriRoute),
	}

}

func (self *UriRoute) RemoveSrv(srv *SrvInfo) error {
	for _, uri := range srv.Uris {
		var sigR *SigUriRoute
		if r, ok := self.router[uri]; ok {
			sigR = r
		} else {
			continue
		}

		if e := sigR.RemoveSrvNode(srv); e != nil {
			return e
		}
	}
	return nil
}

func (self *UriRoute) HandleSrv(srv *SrvInfo) error {
	for _, uri := range srv.Uris {
		var sigR *SigUriRoute
		if r, ok := self.router[uri]; ok {
			sigR = r
		} else {
			sigR = NewSigUriRoute()
			self.router[uri] = sigR
			// srv.r = append(srv.r, sigR)
		}

		if e := sigR.HandleSrvNode(srv); e != nil {
			return e
		}
	}
	return nil
}

func (self *UriRoute) SrvNode(srv string, uri string) ([]*SrvInfo, error) {
	// var *SrvNode
	var sigR *SigUriRoute
	if r, ok := self.router[uri]; ok {
		sigR = r
	}

	return sigR.SrvNode(srv, uri)
	// return nil, nil
}

type SrvFoundRoute struct {
	router map[string]*UriRoute
	// nodes []*SrvNodeInfo
}

func NewSrvFoundRoute() *SrvFoundRoute {
	return &SrvFoundRoute{
		router: make(map[string]*UriRoute),
	}
}

func (self *SrvFoundRoute) HandleSrvNode(node *SrvNodeInfo) error {

	for _, srv := range node.Srv {
		var uriR *UriRoute

		if r, ok := self.router[srv.SrvName]; ok {
			uriR = r
		} else {
			uriR = NewUriRoute()
			self.router[srv.SrvName] = uriR
		}

		if e := uriR.HandleSrv(srv); e != nil {
			return e
		}
	}

	return nil
}

func (self *SrvFoundRoute) RemoveSrvNode(node *SrvNodeInfo) error {
	for _, srv := range node.Srv {
		var uriR *UriRoute

		if r, ok := self.router[srv.SrvName]; ok {
			uriR = r
		} else {
			continue
		}

		if e := uriR.RemoveSrv(srv); e != nil {
			return e
		}
	}
	return nil
}

func (self *SrvFoundRoute) SrvNode(srvName string, uri string) ([]*SrvInfo, error) {
	// var uriR *UriRoute

	if uriR, ok := self.router[srvName]; ok {
		// fmt.Println("uriR: ", uriR)
		// fmt.Printf("srvName: %s, uri: %s\n", srvName, uri)
		return uriR.SrvNode(srvName, uri)
	}

	for k, v := range self.router {
		fmt.Sprintf("k: %s, v: %+v, srvName: %s, uri: %s", k, v, srvName, uri)
		// fmt.Println("k: ", k, "===", "v:", v)
	}
	fmt.Println("router error! ", srvName, "===", uri, "router: ", self.router)
	return nil, errors.New("Not Found Such Srv!")
}
