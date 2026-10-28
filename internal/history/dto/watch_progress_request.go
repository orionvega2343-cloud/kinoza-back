package dto

type SaveProgressRequest struct {
	PositionSec *int `json:"position_sec" binding:"required,gte=0"`
}
