package srv_found

import (
// "fmt"
)

// srv
type SrvInfo struct {
	UUID    string `json:"node_uuid"`
	SrvName string `json:"srv"`
	Addr    string `json:"addr"`
	// Ip      string `json:"ip"`
	// Port    int    `json:"port"`
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

func (self *SrvNodeInfo) SetSrv(srvName string, addr string) {
	// if ip == "" {
	// 	ip = self.Ip
	// }

	srv := &SrvInfo{
		UUID:    self.UUID,
		SrvName: srvName,
		Addr:    addr,
		// Port:    port,
	}

	self.Srv[srvName] = srv

}
