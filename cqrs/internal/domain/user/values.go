package user

import (
	"errors"
	"regexp"
	"unicode/utf8"

	"github.com/google/uuid"
)

type UserID struct {
	value uuid.UUID
}

type HashPassword struct {
	value string
}

type Name struct {
	value string
}

type Image struct {
	value string
}

type Email struct {
	value string
}

type Profile struct {
	value string
}

func NewUserID() (UserID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return UserID{}, err
	}
	return UserID{value: id}, nil
}

func NewHashPassword(value string) HashPassword {
	return HashPassword{value: value}
}

func NewName(value string) (Name, error) {
	length := utf8.RuneCountInString(value)
	if length < 1 || length > 50 {
		return Name{}, errors.New("name must be between 1 and 50 characters")
	}
	return Name{value: value}, nil
}

func NewImage(value string) Image {
	return Image{value: value}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func NewEmail(value string) (Email, error) {
	if !emailRegex.MatchString(value) {
		return Email{value: ""}, errors.New("invalid email format")
	}
	return Email{value: value}, nil
}

func NewProfile(value string) Profile {
	return Profile{value: value}
}

func NewUserIDFromString(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, err
	}
	return UserID{value: id}, nil
}

func (u UserID) String() string {
	return u.value.String()
}

func (h HashPassword) String() string {
	return h.value
}

func (n Name) String() string {
	return n.value
}

func (i Image) String() string {
	return i.value
}

func (e Email) String() string {
	return e.value
}

func (p Profile) String() string {
	return p.value
}
