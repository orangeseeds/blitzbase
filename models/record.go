package model

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
)

type Record struct {
	BaseModel

	collection *Collection
	data       map[string]any
	// remaining expanded relations
}

func NewRecord(c *Collection) *Record {
	return &Record{
		collection: c,
		data:       make(map[string]any),
	}
}

func (r *Record) LoadNullStringMap(data dbx.NullStringMap) *Record {
	for key, val := range data {
		if val.Valid {
			r.Set(key, val.String)
		}
	}
	return r
}

func (r *Record) Load(data map[string]any) {
	for key, val := range data {
		r.Set(key, val)
	}
}

func (r *Record) TableName() string {
	return r.Collection().GetName()
}

func (r *Record) Collection() *Collection {
	return r.collection
}

func (r *Record) Set(key string, val any) {
	switch key {
	case FieldId:
		r.SetID(cast.ToString(val))
	case FieldCreatedAt:
		r.CreatedAt = cast.ToString(val)
	case FieldUpdatedAt:
		r.UpdatedAt = cast.ToString(val)
	default:
		if r.Collection().Schema.HasField(key) {
			r.data[key] = val
		}
	}
}

func (r *Record) Get(key string) any {
	switch CommonFieldName(key) {
	case FieldId:
		return r.GetID()
	default:
		if r.Collection().Schema.HasField(string(key)) {
			v, ok := r.data[string(key)]
			if !ok {
				return nil
			}
			// log.Println(v)
			return v
		}
		return nil
	}
}

func (r *Record) Export() map[string]any {
	toExport := make(map[string]any)
	if r.data != nil {
		toExport = r.data
	}

	log.Println(r)
	toExport[FieldId] = r.Id
	toExport[FieldCreatedAt] = r.CreatedAt
	toExport[FieldUpdatedAt] = r.UpdatedAt
	return toExport
}

func (r Record) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Export())
}

func (r *Record) UnmarshalJSON(data []byte) error {
	res := map[string]any{}

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	r.Load(res)

	return nil
}

func (a *Record) ValidatePassword(password string) bool {
	bytePassword := []byte(password)
	bytePasswordHash := []byte(a.GetString(FieldPassword))

	// comparing the password with the hash
	err := bcrypt.CompareHashAndPassword(bytePasswordHash, bytePassword)

	// nil means it is a match
	return err == nil
}

func (a *Record) SetPassword(password string) error {

	if password == "" {
		return errors.New("The provided plain password is empty")
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	a.Set(FieldPassword, string(hashedPassword))

	a.RefreshToken()
	return nil
}

func (a *Record) RefreshToken() {
	a.Set(string(FieldToken), uuid.NewString())
}

func (m *Record) GetBool(key string) bool {
	return cast.ToBool(m.Get(key))
}

func (m *Record) GetString(key string) string {
	// log.Println("Key", key, m.Get(key))
	return cast.ToString(m.Get(key))
}

func (m *Record) GetInt(key string) int {
	return cast.ToInt(m.Get(key))
}

func (m *Record) GetFloat(key string) float64 {
	return cast.ToFloat64(m.Get(key))
}

func (m *Record) GetTime(key string) time.Time {
	return cast.ToTime(m.Get(key))
}
