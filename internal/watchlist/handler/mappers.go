package handler

import (
	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/dto"
)

// toCreateWatchlistResponse converts a WatchlistItem into a response DTO.
func toCreateWatchlistResponse(item *domain.WatchlistItem) dto.CreateWatchlistResponse {
	return dto.CreateWatchlistResponse{
		ID:      item.ID,
		TitleID: item.TitleID,
		Status:  item.Status,
		AddedAt: item.AddedAt,
	}
}
