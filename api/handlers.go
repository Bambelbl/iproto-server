package api

import (
	"github.com/Bambelbl/iproto-server/packet/request_packet"
	"github.com/Bambelbl/iproto-server/storage"
)

// ADM_STORAGE_SWITCH_READONLY Переводит сторадж в состояние READ_ONLY
func ADM_STORAGE_SWITCH_READONLY(stor *storage.Storage) {
	(*stor).SetState(storage.READ_ONLY)
}

// ADM_STORAGE_SWITCH_READWRITE Переводит сторадж в состояние READ_WRITE
func ADM_STORAGE_SWITCH_READWRITE(stor *storage.Storage) {
	(*stor).SetState(storage.READ_WRITE)
}

// ADM_STORAGE_SWITCH_MAINTENANCE Переводит сторадж в состояние MAINTENANCE
func ADM_STORAGE_SWITCH_MAINTENANCE(stor *storage.Storage) {
	(*stor).SetState(storage.MAINTENANCE)
}

// STORAGE_REPLACE Записывает в сторадж строку по индексу
func STORAGE_REPLACE(stor *storage.Storage, idx int, str string) error {
	return (*stor).SetValue(idx, str)
}

// STORAGE_READ возвращает строку из стораджа по индексу
func STORAGE_READ(stor *storage.Storage, idx int) (string, error) {
	return (*stor).GetValue(idx)
}

func Handler(packet request_packet.IprotoPacketRequest, storage *storage.Storage) (string, uint32) {
	switch packet.Header.Func_id {
	case 65537:
		ADM_STORAGE_SWITCH_READONLY(storage)
		return "", 0
	case 65538:
		ADM_STORAGE_SWITCH_READWRITE(storage)
		return "", 0
	case 65539:
		ADM_STORAGE_SWITCH_MAINTENANCE(storage)
		return "", 0
	case 131073:
		err := STORAGE_REPLACE(storage, packet.Body.Idx, packet.Body.Str)
		if err != nil {
			return err.Error(), 1
		}
		return "", 0
	case 131074:
		body, err := STORAGE_READ(storage, packet.Body.Idx)
		if err != nil {
			return err.Error(), 1
		}
		return body, 0
	default:
		return "Incorrect func_id", 1
	}
}
