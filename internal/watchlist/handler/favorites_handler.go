package handler

import (
	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/dto"
	"kinoza-back/internal/watchlist/service"
	"net/http"

	"log/slog"
	"strconv"

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

func (h *WatchlistHandler) Update(c *gin.Context) {
	var req = dto.UpdateWatchlistRequest{}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		slog.Error("failed to parse watchlist id", "error", err)
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind json", "error", err)
		return
	}
	item := &domain.WatchlistItem{
		ID:     uint(id),
		Status: req.Status,
	}
	if err := h.service.Update(c.Request.Context(), item); err != nil {
		slog.Error("failed to update watchlist item", "error", err)
		return
	}
}

func (h *WatchlistHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		slog.Error("failed to parse watchlist id ", "error", err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		slog.Error("failed to delete watchlist item", "error", err)
		return
	}
}

func (h *WatchlistHandler) List(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 64)

	if err != nil {
		slog.Error("failed to list watchlist items", "error", err)
		return
	}

	items, err := h.service.List(c.Request.Context(), uint(userID))

	if err != nil {
		slog.Error("failed to list watchlist items", "error", err)
		return
	}

	c.JSON(http.StatusOK, items)
}
