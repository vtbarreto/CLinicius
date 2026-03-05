package domain

// ⚠ Architectural violation: domain layer must not depend on infra or repository.
// CLinicius will catch this import.
import (
	"myapp/internal/infra"
	"myapp/internal/repository"
)

// UserService contains business logic for users.
type UserService struct {
	db   *infra.DB
	repo *repository.UserRepository
}

// NewUserService creates a UserService.
// Violation: domain is depending on infrastructure and repository directly.
func NewUserService() *UserService {
	db := infra.Connect("postgres://localhost/myapp")
	return &UserService{
		db:   db,
		repo: &repository.UserRepository{},
	}
}
