package dto

type SignUpRequest struct {
	FirstName string  `json:"firstname" binding:"required"`
	LastName  string  `json:"lastname" binding:"required"`
	Username  string  `json:"username" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Phone     *string `json:"phone"`
	Password  string  `json:"password" binding:"required,min=8"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio"`
	TypeID    *int    `json:"type_id"`
}

type SignUpResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type SignInRequest struct {
	Identifier string `json:"identifier" binding:"required"` // username or email
	Password   string `json:"password" binding:"required"`
}

type SignInResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int64   `json:"id"`
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Balance  float64 `json:"balance"`
		TypeID   int     `json:"type_id"`
	} `json:"user"`
}
