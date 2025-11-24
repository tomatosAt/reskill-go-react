package dto

type UserPasswordPayload struct {
	Username string `json:"username" validate:"required,max=50,regexp=^[a-zA-Z0-9_@ .-]+$"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"omitempty,email,max=100"`
	Remember bool   `json:"remember"`
}
