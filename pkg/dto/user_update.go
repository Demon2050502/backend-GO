package dto

// UpdateUserRequest содержит поля, которые клиент может изменить.
// Используем указатели, чтобы отличать "не указан" от "указан пустым".
type UpdateUserRequest struct {
	FirstName *string `json:"firstname"`
	LastName  *string `json:"lastname"`
	Username  *string `json:"username"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio"`
	TypeID    *int    `json:"type_id"`
	IsActive  *bool   `json:"is_active"`
}

// UserResponse — сокращённый ответ с основными полями.
type UserResponse struct {
	ID        int64   `json:"id"`
	FirstName string  `json:"firstname"`
	LastName  string  `json:"lastname"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Phone     *string `json:"phone,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	TypeID    int     `json:"type_id"`
	IsActive  bool    `json:"is_active"`
	Rating    float64 `json:"rating"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
