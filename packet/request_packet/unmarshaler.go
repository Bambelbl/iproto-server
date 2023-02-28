package request_packet

import (
	"encoding/binary"
	"errors"
	"github.com/vmihailenco/msgpack"
)

// bytes2FuncID from []byte to uint32
func bytes2FuncID(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

// bytes2BodyLength from []byte to uint32
func bytes2BodyLength(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

// bytes2RequestID from []byte to uint32
func bytes2RequestID(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

// bytes2Body from []byte to IprotoBody
func bytes2Body(func_id uint32, body_length uint32, data []byte) (body IprotoBody, err error) {
	if func_id == 0x00020001 {
		buf := make([]byte, body_length)
		err = msgpack.Unmarshal(data, &buf)
		if err != nil {
			return
		}
		body.Idx = int(binary.LittleEndian.Uint32(buf[:4]))
		body.Str = string(buf[4:])
	} else if func_id == 0x00020002 {
		buf := make([]byte, body_length)
		err = msgpack.Unmarshal(data, &buf)
		if err != nil {
			return
		}
		body.Idx = int(binary.LittleEndian.Uint32(buf[:4]))
	}
	return body, nil
}

// Unmarshal from []byte to IprotoPacketRequest
func Unmarshal(data []byte) (requestPacket IprotoPacketRequest, err error) {
	requestPacket.Header.Func_id = bytes2FuncID(data[:4])
	requestPacket.Header.Body_length = bytes2BodyLength(data[4:8])
	requestPacket.Header.Request_id = bytes2RequestID(data[8:12])
	if requestPacket.Header.Body_length > 260 {
		err = errors.New("max length of string is 256 bytes")
	} else {
		requestPacket.Body, err = bytes2Body(requestPacket.Header.Func_id, requestPacket.Header.Body_length, data[12:])
	}
	return
}
