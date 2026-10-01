package domain

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUserById(c *gin.Context)
	UpdateUser(c *gin.Context)
	Refresh(c *gin.Context)
	GetTokenList(c *gin.Context)
	RevokeToken(c *gin.Context)
	RevokeAllTokens(c *gin.Context)
}
