package response_packet

import (
	"encoding/binary"
	"github.com/vmihailenco/msgpack"
)

func ReturnCode2Bytes(data uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, data)
	return res
}

func FuncID2Bytes(data uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, data)
	return res
}

func BodyLength2Bytes(data uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, data)
	return res
}

func RequestID2Bytes(data uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, data)
	return res
}

func Body2Bytes(data string) []byte {
	res, err := msgpack.Marshal(data)
	if err != nil {
		return []byte{}
	}
	return res
}

func Marshal(packet IprotoPacketResponse) []byte {
	var data []byte
	data = append(data, FuncID2Bytes(packet.Header.Func_id)...)
	bodyBytes := Body2Bytes(packet.Body)
	data = append(data, BodyLength2Bytes(uint32(len(bodyBytes)))...)
	data = append(data, RequestID2Bytes(packet.Header.Request_id)...)
	data = append(data, ReturnCode2Bytes(packet.Return_code)...)
	data = append(data, bodyBytes...)
	return data
}
