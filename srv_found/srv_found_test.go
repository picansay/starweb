package srv_found

import (
	"testing"
)

func Test_NewSrvFound(t *testing.T) {
	// num := findByPk(1)
	// t.Log(num)

	node := NewSrvFound()

	node.ListenSrv("srv1")
	node.ListenSrv("srv2")

	t.Log(node)
}
