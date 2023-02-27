package storage

import (
	"testing"
)

type TestCase struct {
	Storage *SimpleStorage
	Idx     int
	Val     string
	State   int
	IsError bool
}

func pointer2Storage(stor SimpleStorage) *SimpleStorage {
	return &stor
}

func TestSimpleStorage_GetState(t *testing.T) {
	cases := []TestCase{
		{
			Storage: pointer2Storage(SimpleStorage{
				state: READ_WRITE,
			}),
			State:   READ_WRITE,
			IsError: false,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				state: READ_ONLY,
			}),
			State:   READ_ONLY,
			IsError: false,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				state: MAINTENANCE,
			}),
			State:   MAINTENANCE,
			IsError: false,
		},
	}
	for caseNum, item := range cases {
		state := item.Storage.GetState()

		if state != item.State {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, state, item.Val)
		}
	}
}

func TestSimpleStorage_GetValue(t *testing.T) {
	data := [1000]string{}
	data[0] = "zero"
	cases := []TestCase{
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: READ_WRITE,
			}),
			Idx:     0,
			Val:     "zero",
			IsError: false,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: READ_ONLY,
			}),
			Idx:     0,
			Val:     "zero",
			IsError: false,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: MAINTENANCE,
			}),
			Idx:     0,
			IsError: true,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: READ_WRITE,
			}),
			Idx:     -1,
			IsError: true,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: READ_WRITE,
			}),
			Idx:     1000,
			IsError: true,
		},
	}
	for caseNum, item := range cases {
		val, err := item.Storage.GetValue(item.Idx)

		if item.IsError && err == nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}

		if !item.IsError && err != nil {
			t.Errorf("[%d] unexpected error: %v", caseNum, err)
		}

		if val != item.Val {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, val, item.Val)
		}
	}
}

func TestSimpleStorage_SetState(t *testing.T) {
	cases := []TestCase{
		{
			Storage: pointer2Storage(SimpleStorage{}),
			State:   READ_WRITE,
			IsError: false,
		},
		{
			Storage: pointer2Storage(SimpleStorage{}),
			State:   READ_ONLY,
			IsError: false,
		},
		{
			Storage: pointer2Storage(SimpleStorage{}),
			State:   MAINTENANCE,
			IsError: false,
		},
	}
	for caseNum, item := range cases {
		item.Storage.SetState(item.State)
		state := item.Storage.GetState()
		if state != item.State {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, state, item.Val)
		}
	}
}

func TestSimpleStorage_SetValue(t *testing.T) {
	data := [1000]string{}
	data[0] = "zero"
	cases := []TestCase{
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: READ_WRITE,
			}),
			Idx:     0,
			Val:     "zero_change",
			IsError: false,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: READ_ONLY,
			}),
			Idx:     0,
			IsError: true,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: MAINTENANCE,
			}),
			Idx:     0,
			IsError: true,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: READ_WRITE,
			}),
			Idx:     -1,
			IsError: true,
		},
		{
			Storage: pointer2Storage(SimpleStorage{
				data:  data,
				state: READ_WRITE,
			}),
			Idx:     1000,
			IsError: true,
		},
	}
	for caseNum, item := range cases {
		err := item.Storage.SetValue(item.Idx, "zero_change")

		if item.IsError && err == nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}

		if !item.IsError && err != nil {
			t.Errorf("[%d] unexpected error: %v", caseNum, err)
		}

		val_change, err := item.Storage.GetValue(item.Idx)
		if err == nil && (item.IsError && val_change == item.Val ||
			!item.IsError && val_change != item.Val) {
			t.Errorf("[%d] wrong results: got %+v, expected %+v",
				caseNum, val_change, item.Val)
		}
	}
}
