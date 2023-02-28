package request_packet

import (
	"encoding/binary"
	"github.com/vmihailenco/msgpack"
	"log"
	"reflect"
	"testing"
)

type TestCase struct {
	Packet  IprotoPacketRequest
	IsError bool
}

func TestUnmarshal(t *testing.T) {
	longText := `В поезде едут 3 юзера и 3 программиста. У юзеров 3 билета, у программистов 1. Заходит контроллер.
	Юзеры показывают билеты, программисты прячутся в туалет. Контроллер стучится в туалет,
	оттуда высовывается рука с билетом. Программисты едут дальше. На обратном пути. У юзеров 1 билет, у программистов ни
	одного. Заходит контроллер.	Юзеры прячутся в туалет. Один из программистов стучит,
	из туалета высовывается рука с билетом. Программисты забирают билет и прячутся в соседний туалет.
	Юзеров ссаживают с поезда.
	Вывод — не всякий алгоритм, доступный программисту, доступен юзеру.`

	cases := []TestCase{
		{
			Packet: IprotoPacketRequest{
				Header: IprotoHeader{
					Func_id:     0x00020001,
					Body_length: 116,
					Request_id:  1,
				},
				Body: IprotoBody{
					Idx: 1,
					Str: "Идет медведь по лесу, видит — машина горит. Сел в нее и сгорел.",
				},
			},
			IsError: false,
		},
		{
			Packet: IprotoPacketRequest{
				Header: IprotoHeader{
					Func_id:     0x00020001,
					Body_length: 1056,
					Request_id:  1,
				},
				Body: IprotoBody{
					Idx: 0,
					Str: longText,
				},
			},
			IsError: true,
		},
		{
			Packet: IprotoPacketRequest{
				Header: IprotoHeader{
					Func_id:     0x00020002,
					Body_length: 4,
					Request_id:  1,
				},
				Body: IprotoBody{
					Idx: 7,
				},
			},
			IsError: false,
		},
	}
	for caseNum, item := range cases {
		input := make([]byte, 12)
		binary.LittleEndian.PutUint32(input[:4], item.Packet.Header.Func_id)
		binary.LittleEndian.PutUint32(input[4:8], item.Packet.Header.Body_length)
		binary.LittleEndian.PutUint32(input[8:12], item.Packet.Header.Request_id)
		if item.Packet.Header.Func_id == 0x00020001 {
			bodyBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bodyBytes[:4], uint32(item.Packet.Body.Idx))
			bodyBytes = append(bodyBytes, []byte(item.Packet.Body.Str)...)
			msgBody, err := msgpack.Marshal(&bodyBytes)
			if err != nil {
				log.Fatalf("Msgpack.marshal error in prepare for test")
			}
			binary.LittleEndian.PutUint32(input[4:8], uint32(len(bodyBytes)))
			input = append(input, msgBody...)
		} else if item.Packet.Header.Func_id == 0x00020002 {
			bodyBytes := make([]byte, 4)
			binary.LittleEndian.PutUint32(bodyBytes[:4], uint32(item.Packet.Body.Idx))
			msgBody, err := msgpack.Marshal(&bodyBytes)
			if err != nil {
				log.Printf("Client: msgpack marshal request error: %s\n", err.Error())
				return
			}
			binary.LittleEndian.PutUint32(input[4:8], uint32(len(bodyBytes)))
			input = append(input, msgBody...)
		}
		packet, err := Unmarshal(input)
		if item.IsError && err == nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}

		if !item.IsError && err != nil {
			t.Errorf("[%d] unexpected error: %v", caseNum, err)
		}

		if err == nil && item.IsError == false && !reflect.DeepEqual(packet, item.Packet) {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, packet, item.Packet)
		}
	}
}
