package response_packet

type IprotoHeader struct {
	Func_id     uint32
	Body_length uint32
	Request_id  uint32
}

type IprotoPacketResponse struct {
	Header      IprotoHeader
	Return_code uint32
	Body        string
}
