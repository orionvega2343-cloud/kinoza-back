package handler

import (
	"kinoza-back/internal/history/domain"
	"kinoza-back/internal/history/dto"
)

func toWatchProgressResponse(p *domain.WatchProgress) dto.WatchProgressResponse {
	return dto.WatchProgressResponse{
		TitleID:     p.TitleID,
		PositionSec: p.PositionSec,
		UpdatedAt:   p.UpdatedAt,
	}
}

func toWatchProgressListResponse(items []domain.WatchProgress) []dto.WatchProgressResponse {
	res := make([]dto.WatchProgressResponse, 0, len(items))

	for _, item := range items {
		res = append(res, toWatchProgressResponse(&item))
	}
	return res
}
