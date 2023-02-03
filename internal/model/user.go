package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID          int
	Login       string
	Password    string
	EncPassword string
}

type Us struct {
	ID                 int
	Login              string
	Password           string
	CategoryOfUser     string
	DateOfRegistration string
	Salt               string
	Algorithm          string
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 8 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncPassword = enc
	}

	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
