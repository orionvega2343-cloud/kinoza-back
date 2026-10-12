package dto

type CreateWatchlistRequest struct {
	Status string `json:"status" binding:"required"`
}
