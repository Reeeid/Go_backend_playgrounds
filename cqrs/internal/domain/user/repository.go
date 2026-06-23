package user

type UserQueryRepository interface {
	GetByID(id UserID) (*User, error)
	GetByEmail(email Email) (*User, error)
}
