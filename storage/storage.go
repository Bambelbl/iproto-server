package storage

type Storage interface {
	GetState() int
	GetValue(idx int) (string, error)
	SetState(state int)
	SetValue(idx int, str string) error
}
