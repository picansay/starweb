package main

import (
	// "errors"
	// "fmt"

	"log"
	// "net"
	"net/http"
	// "time"

	// "github.com/gorilla/mux"
	// "github.com/picansay/starweb"
	"github.com/picansay/starweb"
	sf "github.com/picansay/starweb/srv_found"
	// sn "github.com/picansay/starweb/starnode"
)

// func YourHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(r)
// 	w.Write([]byte("hello!\n"))
// }

func main() {
	// node := sn.NewStarNode()
	// sf "github.com/picansay/starweb/srv_found"
	// go node.Server()
	// node := sf.NewSrvFound()
	// web := starweb.NewStarWeb(node)

	// fmt.Println(web.Addr())
	// // web.HandleFunc("/", YourHandler).Methods("GET")
	// web.HandleFunc("/hello", YourHandler).Methods("GET")
	// log.Fatal(web.Loop())
	node := sf.NewSrvFound()
	// node.ListenSrv("starweb")
	proxy := starweb.NewProxyHandle(node)
	// r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	// r.HandleFunc("/", YourHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", proxy))
}
