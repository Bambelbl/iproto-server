package storage

import (
	"errors"
	"fmt"
	"sync"
)

const (
	SIZE = 1000
)

type SimpleStorage struct {
	state     int
	mutex     sync.RWMutex
	data      [SIZE]string
	dataMutex [SIZE]sync.RWMutex
}

const (
	MAINTENANCE = 0
	READ_ONLY   = 1
	READ_WRITE  = 2
)

func NewSimpleStorageRepo() Storage {
	return &SimpleStorage{state: READ_WRITE}
}

// GetState Return current state of storage
func (s *SimpleStorage) GetState() (state int) {
	s.mutex.RLock()
	state = s.state
	s.mutex.RUnlock()
	return
}

// SetState Set new value of state for storage
func (s *SimpleStorage) SetState(state int) {
	s.mutex.Lock()
	s.state = state
	s.mutex.Unlock()
	return
}

// GetValue Return value from storage by index
func (s *SimpleStorage) GetValue(idx int) (data string, err error) {
	if (*s).GetState() == MAINTENANCE {
		return "", errors.New("storage state doesn't allow this operation")
	}
	if idx < 0 || idx >= SIZE {
		return "", fmt.Errorf("index is out of range: valid index is in [0;%d]", SIZE)
	}
	s.dataMutex[idx].RLock()
	data = s.data[idx]
	s.dataMutex[idx].RUnlock()
	return
}

// SetValue Set value to known index of storage
func (s *SimpleStorage) SetValue(idx int, str string) (err error) {
	if (*s).GetState() != READ_WRITE {
		return errors.New("storage state doesn't allow this operation")
	}
	if idx < 0 || idx >= SIZE {
		return fmt.Errorf("index is out of range: valid index is in [0;%d]", SIZE)
	}
	s.dataMutex[idx].Lock()
	s.data[idx] = str
	s.dataMutex[idx].Unlock()
	return
}
