package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kinoza-back/internal/history/domain"
	"kinoza-back/internal/history/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(svc domain.WatchProgressService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", "user-1")
	})
	h := NewWatchProgressHandler(svc)
	r.PUT("/history/:titleId/progress", h.Save)
	r.GET("/history/continue-watching", h.ListInProgress)
	return r
}

func matchProgress(userID string, titleID, position int) any {
	return mock.MatchedBy(func(p *domain.WatchProgress) bool {
		return p.UserID == userID && p.TitleID == titleID && p.PositionSec == position
	})
}

func TestSave(t *testing.T) {
	tests := []struct {
		name         string
		titleID      string
		body         string
		svcErr       error
		wantStatus   int
		wantCall     bool
		wantPosition int
		wantBody     string
	}{
		{name: "success", titleID: "42", body: `{"position_sec":1200}`, wantStatus: http.StatusOK, wantCall: true, wantPosition: 1200, wantBody: `"position_sec":1200`},
		{name: "zero position", titleID: "42", body: `{"position_sec":0}`, wantStatus: http.StatusOK, wantCall: true, wantPosition: 0},
		{name: "bad title id", titleID: "abc", body: `{"position_sec":1200}`, wantStatus: http.StatusBadRequest},
		{name: "zero title id", titleID: "0", body: `{"position_sec":1200}`, wantStatus: http.StatusBadRequest},
		{name: "bad json", titleID: "42", body: `{bad`, wantStatus: http.StatusBadRequest},
		{name: "missing position", titleID: "42", body: `{}`, wantStatus: http.StatusBadRequest},
		{name: "negative position", titleID: "42", body: `{"position_sec":-5}`, wantStatus: http.StatusBadRequest},
		{name: "invalid position from service", titleID: "42", body: `{"position_sec":1200}`, svcErr: fmt.Errorf("history service: Save, %w", domain.ErrInvalidPosition), wantStatus: http.StatusBadRequest, wantCall: true, wantPosition: 1200, wantBody: `"message":"position must not be negative"`},
		{name: "db error", titleID: "42", body: `{"position_sec":1200}`, svcErr: domain.ErrDatabase, wantStatus: http.StatusInternalServerError, wantCall: true, wantPosition: 1200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.WatchProgressServiceMock)
			if tt.wantCall {
				svc.On("Save", mock.Anything, matchProgress("user-1", 42, tt.wantPosition)).Return(tt.svcErr)
			}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/history/"+tt.titleID+"/progress", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.Contains(t, w.Body.String(), tt.wantBody)
			}
			assert.NotContains(t, w.Body.String(), domain.ErrDatabase.Error())
			svc.AssertExpectations(t)
		})
	}
}

func TestListInProgress(t *testing.T) {
	tests := []struct {
		name       string
		svcErr     error
		wantStatus int
		wantBody   string
	}{
		{name: "success", wantStatus: http.StatusOK, wantBody: "[]"},
		{name: "db error", svcErr: domain.ErrDatabase, wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.WatchProgressServiceMock)
			svc.On("ListInProgress", mock.Anything, "user-1").Return(nil, tt.svcErr)
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/history/continue-watching?user_id=someone-else", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.Equal(t, tt.wantBody, w.Body.String())
			}
			svc.AssertExpectations(t)
		})
	}
}

func TestNoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := new(mocks.WatchProgressServiceMock)
	h := NewWatchProgressHandler(svc)
	r.PUT("/history/:titleId/progress", h.Save)
	r.GET("/history/continue-watching", h.ListInProgress)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "save", method: http.MethodPut, path: "/history/1/progress"},
		{name: "list", method: http.MethodGet, path: "/history/continue-watching"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}

	svc.AssertExpectations(t)
}
