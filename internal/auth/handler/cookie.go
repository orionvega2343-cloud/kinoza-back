package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cookie(c *gin.Context, token string, ttl int) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"refresh_token",
		token,
		ttl,
		"/auth",
		"",
		true,
		true,
	)
}
