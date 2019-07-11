package user

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Model ...
type Model struct {
	Db *sql.DB
}

// Find implement interface
func (m *Model) Find(id int) error {
	return nil
}

// Delete implement interface
func (m *Model) Delete(id int) error {
	return nil
}

// Create implement interface
func (m *Model) Create(u *User) error {
	return nil
}

// FindByLogin impl interface
func (m *Model) FindByLogin(email string, pass string) (User, error) {
	user := User{}

	// Decode password
	err := m.Db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		return User{}, errors.New("User and password not match")
	}

	return user, nil
}
