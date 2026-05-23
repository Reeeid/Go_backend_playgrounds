package command

type CreateUserCommand struct {
	Name  string
	Email string
}

type UpdateUserCommand struct {
	ID    string
	Name  string
	Email string
}

type DeleteUserCommand struct {
	ID string
}

type ChangeUserPasswordCommand struct {
	ID          string
	NewPassword string
}
