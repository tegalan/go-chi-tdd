package user

// User struct
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Store user repository
type Store interface {
	Create(u *User) error
	Delete(id int) error
	Find(id int) error
	FindByLogin(email string, pass string) (User, error)
}
