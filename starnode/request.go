package starnode

type StarNodeRequest struct {
	Client string
	Addr   string
	Uri    string
	Group  string
	Method StarNodeEventType
	Msg    []byte
}
