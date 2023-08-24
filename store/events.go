package store

// type Event interface {
// 	BaseEvent
// }

const (
	INSERT = "insert"
	UPDATE = "update"
	DELETE = "delete"
)

type Event interface {
	FormatSSE() (string, error)
}

type DBHookEvent struct {
	// topic mostly means a collection in sour case
	Message Message
}

func (e *DBHookEvent) FormatSSE() (string, error) {
	return e.Message.FormatSSE()
}
