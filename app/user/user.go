package user

import "golang.org/x/crypto/bcrypt"

// User struct
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// BeforeCreate hook
func (u *User) BeforeCreate() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}
