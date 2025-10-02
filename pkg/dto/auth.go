package dto

type SignUpRequest struct {
	Name      string  `json:"name" binding:"required"`
	Username  string  `json:"username" binding:"required"`
	Email     string  `json:"email" binding:"required,email"`
	Phone     *string `json:"phone"`
	Password  string  `json:"password" binding:"required,min=8"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio"`
	TypeID    *int    `json:"type_id"`
}

type SignUpResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int64   `json:"id"`
		Username string  `json:"username"`
		Email    string  `json:"email"`
	} `json:"user"`
}

type SignInRequest struct {
	Email string `json:"email" binding:"required"`
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
