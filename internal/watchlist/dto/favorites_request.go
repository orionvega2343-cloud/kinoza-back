package dto

type CreateWatchlistRequest struct {
	TitleID uint   `json:"title_id"`
	Status  string `json:"status"`
}

type UpdateWatchlistRequest struct {
	Status string `json:"status"`
}
