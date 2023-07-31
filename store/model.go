package store

import ()

// type CollectionType int
//
// const (
// 	Base CollectionType = 0
// 	Auth CollectionType = 1
// )
//
// type Collection struct {
// 	ID   int
// 	Name string
// 	Type CollectionType
// }

type Model interface {}

type BaseModel struct {
	ID        string
	CreatedAt string
	UpdatedAt string
}

type BaseCollection struct {
    BaseModel

    Name string
    Schema map[string]any
}
