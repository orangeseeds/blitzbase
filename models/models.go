package model

import "github.com/orangeseeds/blitzbase/utils"

type Model interface {
	TableName() string

	// created at updatedat remaining
}

type BaseModel struct {
	Id        string         `db:"id"`
	CreatedAt utils.DateTime `db:"created_at"`
	UpdatedAt utils.DateTime `db:"updated_at"`
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
