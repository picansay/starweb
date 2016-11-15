package srv_found

import (
	"testing"
)

func Test_NewSrvInfo(t *testing.T) {
	// num := findByPk(1)
	// t.Log(num)

	srv := &SrvInfo{
		UUID:    "1111",
		SrvName: "srv1",
		Addr:    "1.1.1.1:80",
		// Port:    80,
	}
	g := NewSrvGroup()
	g.SetSrv(srv)

	sArry, _ := g.GetSrv(srv.SrvName)
	for _, s := range sArry {
		t.Log(s)
	}
	g.RemoveSrv(srv)

	sArry, _ = g.GetSrv(srv.SrvName)
	for _, s := range sArry {
		t.Log(s)
	}
}

func Test_SrvFoundRoute(t *testing.T) {
	// num := findByPk(1)
	// t.Log(num)

	srv := &SrvInfo{
		UUID:    "1111",
		SrvName: "srv1",
		Addr:    "1.1.1.1:80",
		// Port:    80,
	}
	r := NewSrvFoundRoute()
	r.SetSrv(srv)

	sArry, _ := r.GetSrv(srv.SrvName)
	t.Logf("sArry: %+v, len: %d\n", sArry, len(sArry))
	for _, s := range sArry {
		t.Log(s)
	}
	r.RemoveSrv(srv)

	sArry, _ = r.GetSrv(srv.SrvName)
	t.Logf("sArry: %+v, len: %d\n", sArry, len(sArry))
	for _, s := range sArry {
		t.Log(s)
	}
}
