package starweb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sf "github.com/picansay/starweb/srv_found"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "starweb proxy 404 page not found! ", http.StatusNotFound)
}

type ProxyHandle struct {
	node            *sf.SrvFound
	NotFoundHandler http.Handler
}

func NewProxyHandle(node *sf.SrvFound) http.Handler {
	node.ListenSrv("starweb")
	return &ProxyHandle{
		node: node,
	}
}

func (self *ProxyHandle) reqNodes(w http.ResponseWriter, r *http.Request) {
	nodes := self.node.Nodes()
	byt, _ := json.Marshal(nodes)
	w.Write(byt)
}

func (self *ProxyHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	nodeArry, _ := self.node.GetSrv("starweb", r.RequestURI)
	// fmt.Printf("request---> uri: %s\n", r.RequestURI)
	if len(nodeArry) == 0 {
		if self.NotFoundHandler != nil {
			self.NotFoundHandler.ServeHTTP(w, r)
			return
		}
		// self.W
		NotFound(w, r)
		// self.reqNodes(w, r)
		return
	}

	fmt.Println(r)
	srv := nodeArry[0]

	url := fmt.Sprintf("http://%s:%d/%s", srv.Ip, srv.Port, r.RequestURI)
	// c := &http.Client{}
	// r.Host = fmt.Sprintf("%s:%d", srv.Ip, srv.Port)
	// response, err := http.DefaultClient.Do(r)
	// response, _ := c.Do(r)
	response, _ := http.Get(url)
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	w.Write(body)

	// // body, _ := ioutil.ReadAll(response.Body)
	// fmt.Println("node arry: ", nodeArry)
	// byt, _ := json.Marshal(&nodeArry)
	// // byt, _ := json.Marshal(&nodeArry)
	// w.Write(byt)
}
