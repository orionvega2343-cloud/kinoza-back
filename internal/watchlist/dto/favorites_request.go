package dto

type CreateWatchlistRequest struct {
	TitleID int    `json:"title_id" binding:"required,gt=0"`
	Status  string `json:"status" binding:"required"`
}
