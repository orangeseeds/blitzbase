package request

import model "github.com/orangeseeds/blitzbase/models"

type AdminSaveRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

func (r AdminSaveRequest) Model() model.Admin {
	return model.Admin{
		BaseModel: model.BaseModel{},
		Email:     r.Email,
		Password:  r.Password,
	}
}

type AdminAuthWithPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

func (r AdminAuthWithPasswordRequest) Model() model.Admin {
	return model.Admin{}
}

type AdminResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (r AdminResetPasswordRequest) Model() model.Admin {
	return model.Admin{}
}

type AdminConfirmResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	Password        string `json:"password" validate:"required,min=3"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=3,eqfield=Password"`
}

func (r AdminConfirmResetPasswordRequest) Model() model.Admin {
	return model.Admin{}
}

type CollectionSaveRequest struct {
	Name       string               `json:"name" validate:"required,lowercase,ascii"`
	Type       model.CollectionType `json:"type" validate:"required,oneof=base auth"` // base,auth
	Schema     model.Schema         `json:"schema" validate:"required"`
	IndexRule  string               `json:"index_rule"`
	DetailRule string               `json:"detail_rule"`
	CreateRule string               `json:"create_rule"`
	UpdateRule string               `json:"update_rule"`
	DeleteRule string               `json:"delete_rule"`
}

func (r CollectionSaveRequest) Model() model.Collection {
	return model.Collection{
		Name:       r.Name,
		Type:       r.Type,
		Schema:     r.Schema,
		IndexRule:  r.IndexRule,
		DetailRule: r.DetailRule,
		CreateRule: r.CreateRule,
		UpdateRule: r.UpdateRule,
		DeleteRule: r.DeleteRule,
	}
}

type RecordSaveRequest struct {
}

func (r RecordSaveRequest) Model() model.Record {
	return model.Record{}
}

type RecordAuthWithPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

func (r RecordAuthWithPasswordRequest) Model() model.Record {
	return model.Record{}
}

type RecordResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (r RecordResetPasswordRequest) Model() model.Record {
	return model.Record{}
}

type RecordConfirmResetPasswordRequest struct {
	Token           string `json:"token" validate:"required"`
	Password        string `json:"password" validate:"required,min=3"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=3,eqfield=Password"`
}

func (r RecordConfirmResetPasswordRequest) Model() model.Record {
	return model.Record{}
}

type SetSubscriptionRequest struct {
    CollectionId string `json:"collection_id" validate:"required"`
    RecordId     string `json:"record_id" validate:"required"`
    // Type []string `json:"type" validate:oneof=create update delete all` // create, update, delete
}
