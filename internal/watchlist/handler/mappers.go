package handler

import (
	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/dto"
)

func toWatchlistItemResponse(item *domain.WatchlistItem) dto.WatchlistItemResponse {
	return dto.WatchlistItemResponse{
		ID:      item.ID,
		TitleID: item.TitleID,
		Status:  item.Status,
		AddedAt: item.AddedAt,
	}
}

func toWatchlistListResponse(items []domain.WatchlistItem) []dto.WatchlistItemResponse {
	res := make([]dto.WatchlistItemResponse, 0, len(items))

	for _, item := range items {
		res = append(res, toWatchlistItemResponse(&item))
	}
	return res
}
