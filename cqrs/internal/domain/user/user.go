package user

type User struct {
	id      UserID
	name    Name
	image   Image
	email   Email
	profile Profile
}

func NewUser(name Name, image Image, email Email, profile Profile) (*User, error) {
	id, err := NewUserID()
	if err != nil {
		return nil, err
	}
	return &User{
		id:      id,
		name:    name,
		image:   image,
		email:   email,
		profile: profile,
	}, nil
}

func (u *User) Update(name Name, image Image, email Email, profile Profile) {
	u.name = name
	u.image = image
	u.email = email
	u.profile = profile
}

func (u *User) Delete() {
	//command HandlerでUserを削除するためのメソッド
}
