package srv_found

import (
// "fmt"
)

// srv
type SrvInfo struct {
	UUID    string   `json:"node_uuid"`
	SrvName string   `json:"srv"`
	Ip      string   `json:"ip"`
	Port    int      `json:"port"`
	Uris    []string `json:"uris"`
	// r    []*SigUriRoute
}

func (self *SrvInfo) AddUris(uris []string) {
	self.Uris = append(self.Uris, uris...)
}

//node
type SrvNodeInfo struct {
	UUID      string              `json:"uuid"`
	Srv       map[string]*SrvInfo `json:"srv_name"`
	Ip        string              `json:"ip"`
	SrvEnable bool                `json:"srv_enable"`
	Enable    bool                `json:"enable"`
}

func NewSrvNodeInfo() *SrvNodeInfo {
	return &SrvNodeInfo{Srv: make(map[string]*SrvInfo)}
}

func (self *SrvNodeInfo) SetSrv(srvName string, ip string, port int, uris []string) {
	if ip == "" {
		ip = self.Ip
	}

	var srv *SrvInfo

	for _, old_srv := range self.Srv {
		if old_srv.SrvName == srvName {
			srv = old_srv
			break
		}
	}

	if srv == nil {
		srv = &SrvInfo{
			UUID:    self.UUID,
			SrvName: srvName,
			Ip:      ip,
			Port:    port,
			//Uris: uris,
		}

		self.Srv[srvName] = srv
	}
	// fmt.Printf("aaaaaaaaaaa>>>>> node: %+v, uris: %s\n", srv, uris)
	srv.AddUris(uris)
	// fmt.Printf("aaaaaaaaaaa>>>>> node: %+v, uris: %s\n", srv, uris)
}
