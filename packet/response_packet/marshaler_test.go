package response_packet

import (
	"reflect"
	"testing"
)

type TestCase struct {
	Packet IprotoPacketResponse
	Res    []byte
}

func TestMarshal(t *testing.T) {
	cases := []TestCase{
		{
			Packet: IprotoPacketResponse{
				Header: IprotoHeader{
					Func_id:    131074,
					Request_id: 1,
				},
				Return_code: 0,
				Body:        "Идет медведь по лесу, видит — машина горит. Сел в нее и сгорел.",
			},
			Res: []byte{2, 0, 2, 0, 114, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 217, 112, 208, 152, 208, 180, 208, 181, 209,
				130, 32, 208, 188, 208, 181, 208, 180, 208, 178, 208, 181, 208, 180, 209, 140, 32, 208, 191, 208, 190,
				32, 208, 187, 208, 181, 209, 129, 209, 131, 44, 32, 208, 178, 208, 184, 208, 180, 208, 184, 209, 130, 32,
				226, 128, 148, 32, 208, 188, 208, 176, 209, 136, 208, 184, 208, 189, 208, 176, 32, 208, 179, 208, 190,
				209, 128, 208, 184, 209, 130, 46, 32, 208, 161, 208, 181, 208, 187, 32, 208, 178, 32, 208, 189, 208, 181,
				208, 181, 32, 208, 184, 32, 209, 129, 208, 179, 208, 190, 209, 128, 208, 181, 208, 187, 46},
		},
	}
	for caseNum, item := range cases {
		res, err := Marshal(item.Packet)
		if err != nil && !reflect.DeepEqual(res, item.Res) {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, res, item.Res)
		}
	}
}
