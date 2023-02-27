package request_packet

import (
	"encoding/binary"
	"github.com/vmihailenco/msgpack"
)

func bytes2FuncID(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func bytes2BodyLength(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func bytes2RequestID(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func bytes2Body(func_id uint32, data []byte) (body IprotoBody, err error) {
	if func_id == 131073 {
		err = msgpack.Unmarshal(data, body)
		if err != nil {
			return
		}
	} else if func_id == 131074 {
		err = msgpack.Unmarshal(data, body.Idx)
		if err != nil {
			return
		}
		body.Idx = -1
	} else {
		body.Idx = -2
	}
	return body, nil
}

func Unmarshal(data []byte) (requestPacket IprotoPacketRequest, err error) {
	requestPacket.Header.Func_id = bytes2FuncID(data[:4])
	requestPacket.Header.Body_length = bytes2BodyLength(data[4:8])
	requestPacket.Header.Request_id = bytes2RequestID(data[8:12])
	requestPacket.Body, err = bytes2Body(requestPacket.Header.Func_id, data[12:12+requestPacket.Header.Body_length])
	return
}
