package auth

type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required,alphanumspace|email"`
	Password string `json:"password" binding:"required,min=8" validate:"required,min=8"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required,excludesall=0123456789" validate:"required,excludesall=0123456789"`
	LastName  string `json:"last_name" binding:"required,excludesall=0123456789" validate:"required,excludesall=0123456789"`
	Email     string `json:"email" binding:"required,email" validate:"required,email"`
	Username  string `json:"username" binding:"required,alphanumspace" validate:"required,alphanumspace"`
	Password  string `json:"password" binding:"required,min=8" validate:"required,min=8"`
}
