package handler

import (
	"kinoza-back/internal/watchlist/service"
)

type WatchlistHandler struct {
	service *service.WatchlistService
}
// Constructor
func NewWatchlistHandler(service *service.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{
		service: service,
	}
}


