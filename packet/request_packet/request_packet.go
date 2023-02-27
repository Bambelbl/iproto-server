package request_packet

type IprotoBody struct {
	Idx int
	Str string
}

type IprotoHeader struct {
	Func_id     uint32
	Body_length uint32
	Request_id  uint32
}

type IprotoPacketRequest struct {
	Header IprotoHeader
	Body   IprotoBody
}
