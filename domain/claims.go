package domain

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	SocketID        string   `json:"socket_id"`
	FullName        string   `json:"full_name"`
	PhoneNumber     string   `json:"phone_number"`
	Verified        bool     `json:"verified"`
	IsPinAlreadySet bool     `json:"is_pin_already_set "`
	CorporateID     string   `json:"corporate_id"`
	AccessLevel     string   `json:"access_level"`
	Privileges      []string `json:"privileges"`
	jwt.StandardClaims
}

type ClaimsAble interface {
	GetID() string
	GetFullName() string
	GetPhoneNumber() string
	GetVerified() bool
	GetIsPinAlreadySet() bool
	IsLocked() bool

	GetAccessLevel() string
	GetCorporateID() string
	GetPrivileges() []string
}
