package handler

import (
	"net/http"
	"strconv"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/dto"
	"kinoza-back/pkg/response"

	"github.com/gin-gonic/gin"
)

type WatchlistHandler struct {
	service domain.WatchlistService
}

func NewWatchlistHandler(service domain.WatchlistService) *WatchlistHandler {
	return &WatchlistHandler{
		service: service,
	}
}

func (h *WatchlistHandler) Create(c *gin.Context) {
	var req dto.CreateWatchlistRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}
	item := &domain.WatchlistItem{
		TitleID: req.TitleID,
		Status:  req.Status,
		UserID:  c.GetString("user_id"),
	}

	if err := h.service.Create(c.Request.Context(), item); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusCreated, toCreateWatchlistResponse(item))
}

func (h *WatchlistHandler) Update(c *gin.Context) {
	var req dto.UpdateWatchlistRequest

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, response.Error{Message: err.Error(), Code: "INVALID_ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}

	item := &domain.WatchlistItem{
		ID:     id,
		Status: req.Status,
		UserID: c.GetString("user_id"),
	}

	if err := h.service.Update(c.Request.Context(), item); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *WatchlistHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, response.Error{Message: err.Error(), Code: "INVALID_ID"})
		return
	}

	userID := c.GetString("user_id")
	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *WatchlistHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	items, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, items)
}
