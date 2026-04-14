package model

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}