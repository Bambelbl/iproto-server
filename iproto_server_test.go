package main

import (
	"bufio"
	"encoding/binary"
	"github.com/Bambelbl/iproto-server/packet/request_packet"
	"github.com/Bambelbl/iproto-server/packet/response_packet"
	"github.com/vmihailenco/msgpack"
	"log"
	"net"
	"reflect"
	"testing"
)

type TestCase struct {
	input  request_packet.IprotoPacketRequest
	output response_packet.IprotoPacketResponse
}

// Need to start server before testing
func TestServer(t *testing.T) {
	cases := []TestCase{
		{
			input: request_packet.IprotoPacketRequest{
				Header: request_packet.IprotoHeader{
					Func_id:    0x00010001,
					Request_id: 0,
				},
				Body: request_packet.IprotoBody{},
			},
			output: response_packet.IprotoPacketResponse{
				Header: response_packet.IprotoHeader{
					Func_id:     0x00010001,
					Body_length: 0,
					Request_id:  0,
				},
				Return_code: 0,
				Body:        "",
			},
		},
		{
			input: request_packet.IprotoPacketRequest{
				Header: request_packet.IprotoHeader{
					Func_id:    0x00020001,
					Request_id: 1,
				},
				Body: request_packet.IprotoBody{
					Idx: 0,
					Str: "help me pls",
				},
			},
			output: response_packet.IprotoPacketResponse{
				Header: response_packet.IprotoHeader{
					Func_id:     0x00020001,
					Body_length: 44,
					Request_id:  1,
				},
				Return_code: 1,
				Body:        "storage state doesn't allow this operation",
			},
		},
		{
			input: request_packet.IprotoPacketRequest{
				Header: request_packet.IprotoHeader{
					Func_id:    0x00010002,
					Request_id: 2,
				},
				Body: request_packet.IprotoBody{},
			},
			output: response_packet.IprotoPacketResponse{
				Header: response_packet.IprotoHeader{
					Func_id:     0x00010002,
					Body_length: 0,
					Request_id:  2,
				},
				Return_code: 0,
				Body:        "",
			},
		},
		{
			input: request_packet.IprotoPacketRequest{
				Header: request_packet.IprotoHeader{
					Func_id:     0x00020001,
					Body_length: 0,
					Request_id:  3,
				},
				Body: request_packet.IprotoBody{
					Idx: 0,
					Str: "help me pls",
				},
			},
			output: response_packet.IprotoPacketResponse{
				Header: response_packet.IprotoHeader{
					Func_id:     0x00020001,
					Body_length: 0,
					Request_id:  3,
				},
				Return_code: 0,
				Body:        "",
			},
		},
		{
			input: request_packet.IprotoPacketRequest{
				Header: request_packet.IprotoHeader{
					Func_id:     0x00020002,
					Body_length: 0,
					Request_id:  4,
				},
				Body: request_packet.IprotoBody{
					Idx: 0,
				},
			},
			output: response_packet.IprotoPacketResponse{
				Header: response_packet.IprotoHeader{
					Func_id:     0x00020002,
					Body_length: 12,
					Request_id:  4,
				},
				Return_code: 0,
				Body:        "help me pls",
			},
		},
		{
			input: request_packet.IprotoPacketRequest{
				Header: request_packet.IprotoHeader{
					Func_id:    0x00010003,
					Request_id: 5,
				},
				Body: request_packet.IprotoBody{},
			},
			output: response_packet.IprotoPacketResponse{
				Header: response_packet.IprotoHeader{
					Func_id:     0x00010003,
					Body_length: 0,
					Request_id:  5,
				},
				Return_code: 0,
				Body:        "",
			},
		},
		{
			input: request_packet.IprotoPacketRequest{
				Header: request_packet.IprotoHeader{
					Func_id:     0x00020002,
					Body_length: 0,
					Request_id:  6,
				},
				Body: request_packet.IprotoBody{
					Idx: 0,
				},
			},
			output: response_packet.IprotoPacketResponse{
				Header: response_packet.IprotoHeader{
					Func_id:     0x00020002,
					Body_length: 44,
					Request_id:  6,
				},
				Return_code: 1,
				Body:        "storage state doesn't allow this operation",
			},
		},
	}
	for caseNum, item := range cases {
		input := make([]byte, 12)
		binary.LittleEndian.PutUint32(input[:4], item.input.Header.Func_id)
		binary.LittleEndian.PutUint32(input[4:8], item.input.Header.Body_length)
		binary.LittleEndian.PutUint32(input[8:12], item.input.Header.Request_id)
		if item.input.Header.Func_id == 0x00020001 {
			bodyBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bodyBytes[:4], uint32(item.input.Body.Idx))
			bodyBytes = append(bodyBytes, []byte(item.input.Body.Str)...)
			msgBody, err := msgpack.Marshal(&bodyBytes)
			if err != nil {
				log.Fatalf("Msgpack.marshal error in prepare for test")
			}
			binary.LittleEndian.PutUint32(input[4:8], uint32(len(bodyBytes)))
			input = append(input, msgBody...)
		} else if item.input.Header.Func_id == 0x00020002 {
			bodyBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bodyBytes[:4], uint32(item.input.Body.Idx))
			msgBody, err := msgpack.Marshal(&bodyBytes)
			if err != nil {
				log.Printf("Client: msgpack marshal request error: %s\n", err.Error())
				return
			}
			binary.LittleEndian.PutUint32(input[4:8], uint32(len(bodyBytes)))
			input = append(input, msgBody...)
		}
		conn, err := net.Dial("tcp", ":8080")
		defer func(conn net.Conn) {
			err = conn.Close()
			if err != nil {
				log.Printf("Client: connection close error: %s\n", err.Error())
			}
		}(conn)
		if err != nil {
			log.Printf("Client: dial error: %s\n", err.Error())
			return
		}
		_, err = conn.Write(input)
		if err != nil {
			log.Printf("Client: request error: %s\n", err.Error())
			return
		}
		output := make([]byte, 1000)
		_, err = bufio.NewReader(conn).Read(output)
		if err != nil {
			log.Printf("Client: read response error: %s\n", err.Error())
			return
		}
		packet := response_packet.IprotoPacketResponse{
			Header: response_packet.IprotoHeader{
				Func_id:     binary.LittleEndian.Uint32(output[:4]),
				Body_length: binary.LittleEndian.Uint32(output[4:8]),
				Request_id:  binary.LittleEndian.Uint32(output[8:12]),
			},
			Return_code: binary.LittleEndian.Uint32(output[12:16]),
		}
		if packet.Header.Func_id == 0x00020002 ||
			(item.input.Header.Func_id == 0x00020001 && item.output.Return_code == 1) {
			bodyBytes := make([]byte, packet.Header.Body_length)
			err = msgpack.Unmarshal(output[16:16+packet.Header.Body_length], &bodyBytes)
			if err != nil {
				log.Printf("Client: msgpack unmarshal response error: %s\n", err.Error())
				return
			}
			packet.Body = string(bodyBytes)
		}

		if !reflect.DeepEqual(packet, item.output) {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, packet, item.output)
		}

	}
}
