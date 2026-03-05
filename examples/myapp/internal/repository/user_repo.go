package repository

// User represents a stored user record.
type User struct {
	ID   int
	Name string
}

// UserRepository handles user persistence.
type UserRepository struct{}

// FindByID retrieves a user by ID.
func (r *UserRepository) FindByID(id int) *User {
	return &User{ID: id, Name: "example"}
}
