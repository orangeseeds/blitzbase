package model

import (
	"errors"

	"github.com/orangeseeds/blitzbase/utils"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	BaseModel

	Email    string
	Token    string
	Password string // hash

	Rule string // for now all rules
}

func (a *Admin) TableName() string {
	return "_admin_users"
}

func (a *Admin) ValidatePassword(password string) bool {

	bytePassword := []byte(password)
	bytePasswordHash := []byte(a.Password)

	// comparing the password with the hash
	err := bcrypt.CompareHashAndPassword(bytePasswordHash, bytePassword)

	// nil means it is a match
	return err == nil
}
func (a *Admin) SetPassword(password string) error {

	if password == "" {
		return errors.New("The provided plain password is empty")
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	a.Password = string(hashedPassword)

	a.RefreshToken()
	return nil
}

func (a *Admin) RefreshToken() {
	a.Token = utils.RandStr(24)
}

// Provides an export to build the _admin_users table
func (a *Admin) ExportForTable() map[string]string {
	return map[string]string{
		"Id":       string(Text),
		"Email":    string(Text),
		"Token":    string(Text),
		"Password": string(Text),
		"Rule":     string(Text),
	}
}
