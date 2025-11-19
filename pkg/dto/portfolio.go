package dto

type CreatePortfolioRequest struct {
	Token       string                       `json:"token" binding:"required"`
	Title       string                       `json:"title" binding:"required"`
	Description *string                      `json:"description"`
	CategoryID  int                          `json:"category_id" db:"category_id"`
	IsPublic    *bool                        `json:"is_public" db:"is_public"`
	Files       []CreatePortfolioFileRequest `json:"files"` // может быть пустым
}

type CreatePortfolioFileRequest struct {
	FileURL     string  `json:"file_url" db:"file_url" binding:"required"`
	FileType    *string `json:"file_type" db:"file_type"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type CreatePortfolioResponse struct {
	ID int64 `json:"id"`
}

type PortfolioFile struct {
	ID          int64   `json:"id"`
	FileURL     string  `json:"file_url" db:"file_url"`
	FileType    *string `json:"file_type" db:"file_type"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type PortfolioItem struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description *string         `json:"description"`
	CategoryID  int             `json:"category_id" db:"category_id"`
	IsPublic    bool            `json:"is_public" db:"is_public"`
	Files       []PortfolioFile `json:"files"`
}

type GetPortfoliosResponse struct {
	Portfolios []PortfolioItem `json:"portfolios"`
}

type UpdatePortfolioRequest struct {
	Token       string                 `json:"token" binding:"required"`
	PortfolioID int                    `json:"portfolio_id"`
	Title       *string                `json:"title"`
	Description *string                `json:"description"`
	CategoryID  *int                   `json:"category_id" db:"category_id"`
	IsPublic    *bool                  `json:"is_public" db:"is_public"`
	Files       *[]UpdatePortfolioFile `json:"files"`
}

type UpdatePortfolioFile struct {
	FileURL     string  `json:"file_url" db:"file_url" binding:"required"`
	FileType    *string `json:"file_type" db:"file_type"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
