package dto

type GetProfileResponse struct {
	FirstName *string `db:"firstname" json:"firstname,omitempty"`
	LastName  *string `db:"lastname" json:"lastname,omitempty"`
	Username  string  `db:"username" json:"username"`
	Email     *string `db:"email" json:"email,omitempty"`
	Phone     *string `db:"phone" json:"phone,omitempty"`
	AvatarURL *string `db:"avatar_url" json:"avatar_url,omitempty"`
	Bio       *string `db:"bio" json:"bio,omitempty"`
	TypeID    *int64  `db:"type_id" json:"type_id,omitempty"`
	IsActive  *bool   `db:"is_active" json:"is_active,omitempty"`
	CreatedAt *string `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt *string `db:"updated_at" json:"updated_at,omitempty"`
}

// UpdateUserRequest содержит поля, которые клиент может изменить.
type ChangUpdateUserRequest struct {
	Token     string  `json:"token"`
	FirstName *string `json:"firstname"`
	LastName  *string `json:"lastname"`
	Username  *string `json:"username"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio"`
	TypeID    *int    `json:"type_id"`
}

// UserResponse — сокращённый ответ с основными полями.
type ChangUserResponse struct {
	Token     string  `json:"token"`
	FirstName *string `json:"firstname"`
	LastName  *string `json:"lastname"`
	Username  string  `json:"username"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone,omitempty"`
	AvatarURL *string `db:"avatar_url" json:"avatar_url,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	TypeID    *int    `db:"type_id" json:"type_id"`
}
