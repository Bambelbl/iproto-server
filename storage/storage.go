package storage

type Storage interface {

	// GetState Return current state of storage
	GetState() int

	// GetValue Return value from storage by index
	GetValue(idx int) (string, error)

	// SetState Set new value of state for storage
	SetState(state int)

	// SetValue Set value to known index of storage
	SetValue(idx int, str string) error
}
