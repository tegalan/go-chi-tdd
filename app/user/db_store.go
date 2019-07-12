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

// DatabaseStore implement User Store Interface
type DatabaseStore struct {
	DB *sql.DB
}

// Find implement interface
func (m *DatabaseStore) Find(id int) error {
	return nil
}

// Delete implement interface
func (m *DatabaseStore) Delete(id int) error {
	return nil
}

// Create implement interface
func (m *DatabaseStore) Create(u *User) error {
	return nil
}

// FindByLogin impl interface
func (m *DatabaseStore) FindByLogin(email string, pass string) (User, error) {
	user := User{}

	// Decode password
	err := m.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		return User{}, errors.New("User and password not match")
	}

	return user, nil
}
