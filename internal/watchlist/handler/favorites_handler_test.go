package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kinoza-back/internal/watchlist/domain"

	"github.com/gin-gonic/gin"
)

type FakeService struct {
	err       error
	gotUserID string
}

func (f *FakeService) Create(ctx context.Context, item *domain.WatchlistItem) error {
	f.gotUserID = item.UserID
	return f.err
}

func (f *FakeService) Update(ctx context.Context, item *domain.WatchlistItem) error {
	f.gotUserID = item.UserID
	return f.err
}

func (f *FakeService) Delete(ctx context.Context, id int, userID string) error {
	f.gotUserID = userID
	return f.err
}

func (f *FakeService) List(ctx context.Context, userID string) ([]domain.WatchlistItem, error) {
	f.gotUserID = userID
	return nil, f.err
}

var _ domain.WatchlistService = (*FakeService)(nil)

var errDB = errors.New("db error")

func setupRouter(svc domain.WatchlistService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		c.Set("user_id", "user-1")
	})
	h := NewWatchlistHandler(svc)
	r.POST("/watchlist", h.Create)
	r.GET("/watchlist", h.List)
	r.PATCH("/watchlist/:id", h.Update)
	r.DELETE("/watchlist/:id", h.Delete)
	return r
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		svcErr     error
		wantStatus int
	}{
		{name: "success", body: `{"title_id":1,"status":"planned"}`, wantStatus: http.StatusCreated},
		{name: "bad json", body: `{bad`, wantStatus: http.StatusBadRequest},
		{name: "invalid status", body: `{"title_id":1,"status":"planned"}`, svcErr: domain.ErrInvalidStatus, wantStatus: http.StatusBadRequest},
		{name: "already exists", body: `{"title_id":1,"status":"planned"}`, svcErr: domain.ErrAlreadyExists, wantStatus: http.StatusConflict},
		{name: "db error", body: `{"title_id":1,"status":"planned"}`, svcErr: errDB, wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeService{err: tt.svcErr}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/watchlist", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			// на кривой id/json хендлер до сервиса не доходит
			if tt.wantStatus != http.StatusBadRequest || tt.svcErr != nil {
				if svc.gotUserID != "user-1" {
					t.Errorf("expected user_id user-1, got %q", svc.gotUserID)
				}
			}
			if tt.wantStatus == http.StatusInternalServerError && strings.Contains(w.Body.String(), errDB.Error()) {
				t.Errorf("db error text leaked to client: %s", w.Body.String())
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		body       string
		svcErr     error
		wantStatus int
	}{
		{name: "success", id: "1", body: `{"status":"watched"}`, wantStatus: http.StatusOK},
		{name: "bad id", id: "abc", body: `{"status":"watched"}`, wantStatus: http.StatusBadRequest},
		{name: "bad json", id: "1", body: `{bad`, wantStatus: http.StatusBadRequest},
		{name: "not found", id: "1", body: `{"status":"watched"}`, svcErr: domain.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "invalid status", id: "1", body: `{"status":"watched"}`, svcErr: domain.ErrInvalidStatus, wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeService{err: tt.svcErr}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/watchlist/"+tt.id, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			// на кривой id/json хендлер до сервиса не доходит
			if tt.wantStatus != http.StatusBadRequest || tt.svcErr != nil {
				if svc.gotUserID != "user-1" {
					t.Errorf("expected user_id user-1, got %q", svc.gotUserID)
				}
			}
			if tt.wantStatus == http.StatusInternalServerError && strings.Contains(w.Body.String(), errDB.Error()) {
				t.Errorf("db error text leaked to client: %s", w.Body.String())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		svcErr     error
		wantStatus int
	}{
		{name: "success", id: "1", wantStatus: http.StatusOK},
		{name: "bad id", id: "abc", wantStatus: http.StatusBadRequest},
		{name: "not found", id: "1", svcErr: domain.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "db error", id: "1", svcErr: errDB, wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &FakeService{err: tt.svcErr}
			r := setupRouter(svc)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/watchlist/"+tt.id, nil)
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			// на кривой id/json хендлер до сервиса не доходит
			if tt.wantStatus != http.StatusBadRequest || tt.svcErr != nil {
				if svc.gotUserID != "user-1" {
					t.Errorf("expected user_id user-1, got %q", svc.gotUserID)
				}
			}
			if tt.wantStatus == http.StatusInternalServerError && strings.Contains(w.Body.String(), errDB.Error()) {
				t.Errorf("db error text leaked to client: %s", w.Body.String())
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
		{name: "db error", svcErr: errDB, wantStatus: http.StatusInternalServerError},
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
