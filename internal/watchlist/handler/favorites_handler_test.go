package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kinoza-back/internal/watchlist/domain"

	"github.com/gin-gonic/gin"
)

type FakeService struct {
	err        error
	gotUserID  string
	gotTitleID int
}

func (f *FakeService) Create(ctx context.Context, item *domain.WatchlistItem) error {
	f.gotUserID = item.UserID
	f.gotTitleID = item.TitleID
	return f.err
}

func (f *FakeService) Update(ctx context.Context, item *domain.WatchlistItem) error {
	f.gotUserID = item.UserID
	f.gotTitleID = item.TitleID
	return f.err
}

func (f *FakeService) Delete(ctx context.Context, titleID int, userID string) error {
	f.gotUserID = userID
	f.gotTitleID = titleID
	return f.err
}

func (f *FakeService) List(ctx context.Context, userID string) ([]domain.WatchlistItem, error) {
	f.gotUserID = userID
	return nil, f.err
}

var _ domain.WatchlistService = (*FakeService)(nil)

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
			svc := &FakeService{err: tt.svcErr}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/watchlist/"+tt.titleID, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			if svc.gotUserID != tt.wantUserID {
				t.Errorf("expected user_id %q, got %q", tt.wantUserID, svc.gotUserID)
			}
			if svc.gotTitleID != tt.wantTitleID {
				t.Errorf("expected title_id %d, got %d", tt.wantTitleID, svc.gotTitleID)
			}
			if strings.Contains(w.Body.String(), domain.ErrDB.Error()) {
				t.Errorf("db error leaked: %s", w.Body.String())
			}
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
			svc := &FakeService{err: tt.svcErr}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/watchlist/"+tt.titleID, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			if svc.gotUserID != tt.wantUserID {
				t.Errorf("expected user_id %q, got %q", tt.wantUserID, svc.gotUserID)
			}
			if svc.gotTitleID != tt.wantTitleID {
				t.Errorf("expected title_id %d, got %d", tt.wantTitleID, svc.gotTitleID)
			}
			if strings.Contains(w.Body.String(), domain.ErrDB.Error()) {
				t.Errorf("db error leaked: %s", w.Body.String())
			}
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
			svc := &FakeService{err: tt.svcErr}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/watchlist/"+tt.titleID, nil)
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			if svc.gotUserID != tt.wantUserID {
				t.Errorf("expected user_id %q, got %q", tt.wantUserID, svc.gotUserID)
			}
			if svc.gotTitleID != tt.wantTitleID {
				t.Errorf("expected title_id %d, got %d", tt.wantTitleID, svc.gotTitleID)
			}
			if strings.Contains(w.Body.String(), domain.ErrDB.Error()) {
				t.Errorf("db error leaked: %s", w.Body.String())
			}
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
			svc := &FakeService{err: tt.svcErr}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/watchlist?user_id=someone-else", nil)
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			if svc.gotUserID != "user-1" {
				t.Errorf("expected user_id user-1, got %q", svc.gotUserID)
			}
		})
	}
}

func TestNoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewWatchlistHandler(&FakeService{})
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

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
			}
		})
	}
}
