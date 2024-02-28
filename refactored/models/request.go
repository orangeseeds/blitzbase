package model

type Request struct {
	BaseModel

	URL      string
	Method   string
	Auth     string // guest, auth
	RemoteIP string
}

func (r *Request) TableName() string {
	return "_requests"
}

