package dto

type CreateWatchlistRequest struct {
	TitleID int    `json:"title_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}
