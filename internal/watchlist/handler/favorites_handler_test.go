package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kinoza-back/internal/watchlist/domain"
	"kinoza-back/internal/watchlist/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(svc domain.WatchlistService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", "user-1")
	})
	h := NewWatchlistHandler(svc)
	r.POST("/watchlist/:titleId", h.Create)
	r.GET("/watchlist", h.List)
	r.PATCH("/watchlist/:titleId", h.Update)
	r.DELETE("/watchlist/:titleId", h.Delete)
	return r
}

// matchItem проверяет, что хендлер передал в сервис правильные user_id и title_id.
func matchItem(userID string, titleID int) any {
	return mock.MatchedBy(func(item *domain.WatchlistItem) bool {
		return item.UserID == userID && item.TitleID == titleID
	})
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		titleID     string
		body        string
		svcErr      error
		wantStatus  int
		wantUserID  string
		wantTitleID int
	}{
		{name: "success", titleID: "42", body: `{"status":"planned"}`, wantStatus: http.StatusCreated, wantUserID: "user-1", wantTitleID: 42},
		{name: "bad title id", titleID: "abc", body: `{"status":"planned"}`, wantStatus: http.StatusBadRequest},
		{name: "zero title id", titleID: "0", body: `{"status":"planned"}`, wantStatus: http.StatusBadRequest},
		{name: "bad json", titleID: "42", body: `{bad`, wantStatus: http.StatusBadRequest},
		{name: "invalid status", titleID: "42", body: `{"status":"planned"}`, svcErr: domain.ErrInvalidStatus, wantStatus: http.StatusBadRequest, wantUserID: "user-1", wantTitleID: 42},
		{name: "db error", titleID: "42", body: `{"status":"planned"}`, svcErr: domain.ErrDB, wantStatus: http.StatusInternalServerError, wantUserID: "user-1", wantTitleID: 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.WatchlistServiceMock)
			if tt.wantUserID != "" {
				svc.On("Create", mock.Anything, matchItem(tt.wantUserID, tt.wantTitleID)).Return(tt.svcErr)
			}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/watchlist/"+tt.titleID, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.NotContains(t, w.Body.String(), domain.ErrDB.Error())
			svc.AssertExpectations(t)
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name        string
		titleID     string
		body        string
		svcErr      error
		wantStatus  int
		wantUserID  string
		wantTitleID int
	}{
		{name: "success", titleID: "42", body: `{"status":"watched"}`, wantStatus: http.StatusOK, wantUserID: "user-1", wantTitleID: 42},
		{name: "bad title id", titleID: "abc", body: `{"status":"watched"}`, wantStatus: http.StatusBadRequest},
		{name: "zero title id", titleID: "0", body: `{"status":"watched"}`, wantStatus: http.StatusBadRequest},
		{name: "bad json", titleID: "42", body: `{bad`, wantStatus: http.StatusBadRequest},
		{name: "not found", titleID: "42", body: `{"status":"watched"}`, svcErr: domain.ErrNotFound, wantStatus: http.StatusNotFound, wantUserID: "user-1", wantTitleID: 42},
		{name: "invalid status", titleID: "42", body: `{"status":"watched"}`, svcErr: domain.ErrInvalidStatus, wantStatus: http.StatusBadRequest, wantUserID: "user-1", wantTitleID: 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.WatchlistServiceMock)
			if tt.wantUserID != "" {
				svc.On("Update", mock.Anything, matchItem(tt.wantUserID, tt.wantTitleID)).Return(tt.svcErr)
			}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/watchlist/"+tt.titleID, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.NotContains(t, w.Body.String(), domain.ErrDB.Error())
			svc.AssertExpectations(t)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name        string
		titleID     string
		svcErr      error
		wantStatus  int
		wantUserID  string
		wantTitleID int
	}{
		{name: "success", titleID: "42", wantStatus: http.StatusOK, wantUserID: "user-1", wantTitleID: 42},
		{name: "bad title id", titleID: "abc", wantStatus: http.StatusBadRequest},
		{name: "zero title id", titleID: "0", wantStatus: http.StatusBadRequest},
		{name: "not found", titleID: "42", svcErr: domain.ErrNotFound, wantStatus: http.StatusNotFound, wantUserID: "user-1", wantTitleID: 42},
		{name: "db error", titleID: "42", svcErr: domain.ErrDB, wantStatus: http.StatusInternalServerError, wantUserID: "user-1", wantTitleID: 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.WatchlistServiceMock)
			if tt.wantUserID != "" {
				svc.On("Delete", mock.Anything, tt.wantTitleID, tt.wantUserID).Return(tt.svcErr)
			}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/watchlist/"+tt.titleID, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.NotContains(t, w.Body.String(), domain.ErrDB.Error())
			svc.AssertExpectations(t)
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name       string
		svcErr     error
		wantStatus int
	}{
		{name: "success", wantStatus: http.StatusOK},
		{name: "db error", svcErr: domain.ErrDB, wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mocks.WatchlistServiceMock)
			// user_id берётся из контекста, а не из query — мок ждёт "user-1"
			svc.On("List", mock.Anything, "user-1").Return(nil, tt.svcErr)
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/watchlist?user_id=someone-else", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			svc.AssertExpectations(t)
		})
	}
}

func TestNoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := new(mocks.WatchlistServiceMock)
	h := NewWatchlistHandler(svc)
	r.POST("/watchlist/:titleId", h.Create)
	r.GET("/watchlist", h.List)
	r.PATCH("/watchlist/:titleId", h.Update)
	r.DELETE("/watchlist/:titleId", h.Delete)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "create", method: http.MethodPost, path: "/watchlist/1"},
		{name: "list", method: http.MethodGet, path: "/watchlist"},
		{name: "update", method: http.MethodPatch, path: "/watchlist/1"},
		{name: "delete", method: http.MethodDelete, path: "/watchlist/1"},
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
