package store

type CollectionType int

const (
	Base CollectionType = 0
    Auth CollectionType = 1
)

type Collection struct {
	ID   int
	Name string
	Type CollectionType
}
