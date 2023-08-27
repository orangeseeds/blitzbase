package store

type Model interface {
	TableName() string
}

type BaseModel struct {
	Id      string
	Created string
	Updated string
}
