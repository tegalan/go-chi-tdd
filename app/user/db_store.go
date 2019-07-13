package user

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

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

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	var id int
	if err := m.DB.QueryRow("INSERT INTO users (name, email, password) VALUES($1, $2, $3) RETURNING id", u.Name, u.Email, u.Password).Scan(&id); err != nil {
		return err
	}

	u.ID = id
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
