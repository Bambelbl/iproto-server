package request_packet

import (
	"github.com/vmihailenco/msgpack"
	"reflect"
	"testing"
)

type TestCase struct {
	Packet  IprotoPacketRequest
	Input   []byte
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

	longBody, _ := msgpack.Marshal(IprotoBody{Idx: 1,
		Str: longText,
	})

	cases := []TestCase{
		{
			Packet: IprotoPacketRequest{
				Header: IprotoHeader{
					Func_id:     131073,
					Body_length: 132,
					Request_id:  1,
				},
				Body: IprotoBody{Idx: 1,
					Str: "Идет медведь по лесу, видит — машина горит. Сел в нее и сгорел.",
				},
			},
			Input: []byte{1, 0, 2, 0, 132, 0, 0, 0, 1, 0, 0, 0, 130, 163, 73, 100, 120, 211, 0, 0, 0, 0, 0, 0,
				0, 1, 163, 83, 116, 114, 217, 112, 208, 152, 208, 180, 208, 181, 209, 130, 32, 208, 188, 208, 181, 208,
				180, 208, 178, 208, 181, 208, 180, 209, 140, 32, 208, 191, 208, 190, 32, 208, 187, 208, 181, 209, 129,
				209, 131, 44, 32, 208, 178, 208, 184, 208, 180, 208, 184, 209, 130, 32, 226, 128, 148, 32, 208, 188, 208,
				176, 209, 136, 208, 184, 208, 189, 208, 176, 32, 208, 179, 208, 190, 209, 128, 208, 184, 209, 130, 46,
				32, 208, 161, 208, 181, 208, 187, 32, 208, 178, 32, 208, 189, 208, 181, 208, 181, 32, 208, 184, 32, 209,
				129, 208, 179, 208, 190, 209, 128, 208, 181, 208, 187, 46},
			IsError: false,
		},
		{
			Packet: IprotoPacketRequest{
				Header: IprotoHeader{
					Func_id:     131073,
					Body_length: 1073,
					Request_id:  1,
				},
				Body: IprotoBody{Idx: 0,
					Str: "",
				},
			},
			Input:   append([]byte{1, 0, 2, 0, 49, 4, 0, 0, 1, 0, 0, 0}, longBody...),
			IsError: true,
		},
		{
			Packet: IprotoPacketRequest{
				Header: IprotoHeader{
					Func_id:     131074,
					Body_length: 9,
					Request_id:  1,
				},
				Body: IprotoBody{Idx: 7,
					Str: "",
				},
			},
			Input:   []byte{2, 0, 2, 0, 9, 0, 0, 0, 1, 0, 0, 0, 211, 0, 0, 0, 0, 0, 0, 0, 7},
			IsError: false,
		},
	}
	for caseNum, item := range cases {
		packet, err := Unmarshal(item.Input)
		if item.IsError && err == nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}

		if !item.IsError && err != nil {
			t.Errorf("[%d] unexpected error: %v", caseNum, err)
		}

		if !reflect.DeepEqual(packet, item.Packet) {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, packet, item.Packet)
		}
	}
}
