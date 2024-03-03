package model

type Model interface {
	GetID() string
	TableName() string
	SetID(id string)

	// created at updatedat remaining
}

type BaseModel struct {
	Id        string `db:"Id"`
	CreatedAt string `db:"CreatedAt"`
	UpdatedAt string `db:"UpdatedAt"`
}

func NewBaseModel(id string) *BaseModel {
	return &BaseModel{
		Id: id,
	}
}

func (b *BaseModel) GetID() string {
	return b.Id
}
func (b *BaseModel) SetID(id string) {
	b.Id = id
}
