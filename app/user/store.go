package user

// Store user repository
type Store interface {
	Create(u *User) error
	Delete(id int) error
	Find(id int) error
	FindByLogin(email string, pass string) (User, error)
}
