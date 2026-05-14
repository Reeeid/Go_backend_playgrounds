package user

import "time"

type Credential struct {
	userID    UserID
	password  HashPassword
	updatedAt time.Time
	createdAt time.Time
}

func NewCredential(userID UserID, password HashPassword) *Credential {
	return &Credential{
		userID:    userID,
		password:  password,
		updatedAt: time.Now(),
		createdAt: time.Now(),
	}
}

func (c *Credential) Update(password HashPassword) {
	c.password = password
	c.updatedAt = time.Now()
}

func (c *Credential) Delete() {
	//command HandlerでCredentialを削除するためのメソッド
}
