package models

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// SignupRequest is used for the initial email/password submission
type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

// VerifyRequest is used when the user submits the OTP code
type VerifyRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

// ForgotPasswordRequest is used to trigger the OTP email
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest is the final step where the user provides a new password
type ResetPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}


//Rooms
type CreateRoomRequest struct {
	Name    string  `json:"name,omitempty"` // Nullable for 1-on-1
	IsGroup bool    `json:"is_group"`
	Members []int64 `json:"members"`
	UserId  int64   `json:"-"`
}