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

func getUserID(c *gin.Context) (string, bool) {
	userID := c.GetString("user_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, response.Error{Message: "unauthorized", Code: "UNAUTHORIZED"})
		return "", false
	}
	return userID, true
}

func (h *WatchlistHandler) Create(c *gin.Context) {
	var req dto.CreateWatchlistRequest
	userID, ok := getUserID(c)

	if !ok {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error{Message: "invalid id", Code: "FAILED_TO_BIND"})
		return
	}
	item := &domain.WatchlistItem{
		TitleID: req.TitleID,
		Status:  req.Status,
		UserID:  userID,
	}

	if err := h.service.Create(c.Request.Context(), item); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusCreated, toWatchlistItemResponse(item))
}

func (h *WatchlistHandler) Update(c *gin.Context) {
	var req dto.UpdateWatchlistRequest
	userID, ok := getUserID(c)

	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, response.Error{Message: "invalid id", Code: "INVALID_ID"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}

	item := &domain.WatchlistItem{
		ID:     id,
		Status: req.Status,
		UserID: userID,
	}

	if err := h.service.Update(c.Request.Context(), item); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *WatchlistHandler) Delete(c *gin.Context) {
	userID, ok := getUserID(c)

	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(400, response.Error{Message: "invalid id", Code: "INVALID_ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *WatchlistHandler) List(c *gin.Context) {
	userID, ok := getUserID(c)

	if !ok {
		return
	}

	items, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, toWatchlistListResponse(items))
}
