package dto

type RegisterRequest struct {
	Name     string  `json:"name" validate:"required,min=2,max=150"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Phone    *string `json:"phone" validate:"omitempty,max=20"`
	Role     string  `json:"role" validate:"omitempty,oneof=user promotor"`
	Gender   *string `json:"gender" validate:"omitempty,oneof=male female other"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
	Name   *string `json:"name" validate:"omitempty,min=2,max=150"`
	Phone  *string `json:"phone" validate:"omitempty,max=20"`
	Gender *string `json:"gender" validate:"omitempty,oneof=male female other"`
}
