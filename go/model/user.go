package model

import (
	"fmt"
	"golang.org/x/crypto/scrypt"
)

var Users []User

type User struct {
	Login    string
	Password string
	Access   []string
}

func NewUser(l, pw string, access ...string) *User {
	return &User{
		Login:    l,
		Password: pw,
		Access:   access,
	}
}

func (u *User) ValidatePw(pw string) bool {
	return u.Password == encrypt(pw+u.Login)
}

//func (u *User) ChangePw(oldpw, newpw string) bool {
//	if u.ValidatePw(oldpw) {
//		u.Password = encrypt(newpw)
//		return true
//	} else {
//		return false
//	}
//}

func (u *User) Authorize(target string) bool {
	for _, a := range u.Access {
		if a == target {
			return true
		}
	}
	return false
}

func (u *User) String() string {
	return fmt.Sprintf("%s: [%s] %v", u.Login, u.Password, u.Access)
}

func encrypt(pw string) string {
	dk, err := scrypt.Key([]byte(pw), []byte("#K?@1"), 16384, 4, 1, 32)
	if err != nil {
		panic(fmt.Errorf("Encryption error %s \n", err))
	}
	return fmt.Sprintf("%x", dk)
}

func Encrypt(pw string) string {
	return encrypt(pw)
}
