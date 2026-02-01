package response

import "time"

type UserResponse struct {
	ID        uint      `json:"id"`
	EntityID  uint      `json:"entity_id"`
	Username  string    `json:"username"`
	LastLogin time.Time `json:"last_login"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type RegisterResponse struct {
	Entity EntityResponse `json:"entity"`
	User   UserResponse   `json:"user"`
}