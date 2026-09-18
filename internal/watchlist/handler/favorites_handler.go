package handler

import (
	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/dto"
	"kinoza-back/internal/watchlist/service"
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
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

func (h *WatchlistHandler) Create(c *gin.Context) {
	var req = dto.CreateWatchlistRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind json,", "error", err)
		return
	}

	item := &domain.WatchlistItem{
		TitleID: req.TitleID,
		Status:  req.Status,
	}

	if err := h.service.Create(c.Request.Context(), item); err != nil {
		slog.Error("failed to create watchlist item", "error", err)
		return
	}
	c.JSON(http.StatusCreated, dto.CreateWatchlistResponse{
		ID:      item.ID,
		TitleID: item.TitleID,
		Status:  item.Status,
		AddedAt: item.AddedAt,
	})
}

// TODO: make an update func
