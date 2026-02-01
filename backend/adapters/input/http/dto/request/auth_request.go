package request

type RegisterRequest struct {
	Type         string `json:"type" binding:"required"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}