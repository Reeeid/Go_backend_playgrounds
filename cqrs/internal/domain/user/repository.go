package user

type UserCommandRepository interface {
	Save(user *User) error
	Delete(id UserID) error
	Update(user *User) error
}

type UserQueryRepository interface {
	GetByID(id UserID) (*User, error)
	GetByEmail(email Email) (*User, error)
}
