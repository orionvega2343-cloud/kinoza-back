package handler

import (
	"net/http"
	"strconv"

	"kinoza-back/internal/history/domain"
	"kinoza-back/internal/history/dto"
	"kinoza-back/pkg/response"

	"github.com/gin-gonic/gin"
)

type WatchProgressHandler struct {
	service domain.WatchProgressService
}

func NewWatchProgressHandler(service domain.WatchProgressService) *WatchProgressHandler {
	return &WatchProgressHandler{
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

func getTitleID(c *gin.Context) (int, bool) {
	titleID, err := strconv.Atoi(c.Param("titleId"))

	if err != nil || titleID <= 0 {
		c.JSON(http.StatusBadRequest, response.Error{Message: "invalid title id", Code: "INVALID_TITLE_ID"})
		return 0, false
	}
	return titleID, true
}

func (h *WatchProgressHandler) Save(c *gin.Context) {
	var req dto.SaveProgressRequest
	userID, ok := getUserID(c)

	if !ok {
		return
	}
	titleID, ok := getTitleID(c)
	if !ok {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}
	p := &domain.WatchProgress{
		UserID:      userID,
		TitleID:     titleID,
		PositionSec: *req.PositionSec,
	}

	if err := h.service.Save(c.Request.Context(), p); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(http.StatusOK, toWatchProgressResponse(p))
}

func (h *WatchProgressHandler) ListInProgress(c *gin.Context) {
	userID, ok := getUserID(c)

	if !ok {
		return
	}
	items, err := h.service.ListInProgress(c.Request.Context(), userID)
	if err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}

	c.JSON(http.StatusOK, toWatchProgressListResponse(items))
}
