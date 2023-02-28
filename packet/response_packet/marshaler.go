package response_packet

import (
	"encoding/binary"
	"github.com/vmihailenco/msgpack"
)

// ReturnCode2Bytes from uint32 to []byte
func ReturnCode2Bytes(data uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, data)
	return res
}

// FuncID2Bytes from uint32 to []byte
func FuncID2Bytes(data uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, data)
	return res
}

// BodyLength2Bytes from uint32 to []byte
func BodyLength2Bytes(data uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, data)
	return res
}

// RequestID2Bytes from uint32 to []byte
func RequestID2Bytes(data uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, data)
	return res
}

// Body2Bytes from string to []byte
func Body2Bytes(data string) (res []byte, err error) {
	if data != "" {
		res, err = msgpack.Marshal(&data)
		if err != nil {
			return
		}
	}
	return res, nil
}

// Marshal from IprotoPacketResponse to []byte
func Marshal(packet IprotoPacketResponse) (data []byte, err error) {
	data = append(data, FuncID2Bytes(packet.Header.Func_id)...)
	bodyBytes, err := Body2Bytes(packet.Body)
	data = append(data, BodyLength2Bytes(uint32(len(bodyBytes)))...)
	data = append(data, RequestID2Bytes(packet.Header.Request_id)...)
	data = append(data, ReturnCode2Bytes(packet.Return_code)...)
	data = append(data, bodyBytes...)
	return data, err
}
