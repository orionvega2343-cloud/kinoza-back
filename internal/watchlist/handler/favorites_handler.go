package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/dto"
	"kinoza-back/internal/watchlist/service"

	"github.com/gin-gonic/gin"
)

type WatchlistHandler struct {
	service *service.WatchlistService
}

func NewWatchlistHandler(service *service.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{
		service: service,
	}
}

func (h *WatchlistHandler) Create(c *gin.Context) {
	var req dto.CreateWatchlistRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind json", "error", err)
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

	c.JSON(http.StatusCreated, toCreateWatchlistResponse(item))
}

func (h *WatchlistHandler) Update(c *gin.Context) {
	var req dto.UpdateWatchlistRequest

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Error("failed to parse watchlist id", "error", err)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind json", "error", err)
		return
	}

	item := &domain.WatchlistItem{
		ID:     id,
		Status: req.Status,
	}

	if err := h.service.Update(c.Request.Context(), item); err != nil {
		slog.Error("failed to update watchlist item", "error", err)
		return
	}
}

func (h *WatchlistHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Error("failed to parse watchlist id", "error", err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		slog.Error("failed to delete watchlist item", "error", err)
		return
	}
}

func (h *WatchlistHandler) List(c *gin.Context) {
	userID := c.Query("user_id")

	items, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		slog.Error("failed to list watchlist items", "error", err)
		return
	}

	c.JSON(http.StatusOK, items)
}
