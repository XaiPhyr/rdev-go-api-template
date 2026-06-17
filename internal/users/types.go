package users

import "github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"

type UserRequest struct {
	FirstName *string `json:"first_name" binding:"required,excludesall=0123456789" validate:"required,excludesall=0123456789"`
	LastName  *string `json:"last_name" binding:"required,excludesall=0123456789" validate:"required,excludesall=0123456789"`
	Email     *string `json:"email" binding:"required,email" validate:"required,email"`
	Username  *string `json:"username" binding:"required,alphanumspace" validate:"required,alphanumspace"`
	Password  *string `json:"password" binding:"required,min=8" validate:"required,min=8"`
	UserType  *string `json:"user_type" bindng:"required"`
}

type UserResponse struct {
	// struct for response goes here!
}

func (req UserRequest) ToModel(user *models.User) *models.User {
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Username != nil {
		user.Username = *req.Username
	}

	if req.Password != nil {
		user.Password = *req.Password
	}

	if req.UserType != nil {
		user.UserType = *req.UserType
	}

	return user
}
