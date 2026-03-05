package handler

// ⚠ Architectural violation: handler layer must not depend on repository directly.
// CLinicius will catch this import.
import (
	"fmt"
	"myapp/internal/repository"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	repo *repository.UserRepository
}

// NewUserHandler creates a UserHandler.
// Violation: handler is bypassing the domain and importing repository directly.
func NewUserHandler() *UserHandler {
	return &UserHandler{repo: &repository.UserRepository{}}
}

// Handle processes a user request.
func (h *UserHandler) Handle(id int) {
	user := h.repo.FindByID(id)
	fmt.Printf("user: %+v\n", user)
}
