package handler

import (
	"kinoza-back/internal/auth/domain"
	"kinoza-back/internal/auth/dto"
	"kinoza-back/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandlerImpl struct {
	svc       domain.UserService
	cookieTTL int
}

func NewUserHandler(svc domain.UserService, cookieTTL int) *UserHandlerImpl {
	return &UserHandlerImpl{svc: svc, cookieTTL: cookieTTL}
}

func (h *UserHandlerImpl) Register(c *gin.Context) {
	var req dto.UserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}
	ctx := c.Request.Context()
	registered, err := h.svc.Register(ctx, new(toUserDomainUser(req)))
	if err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(201, toUserResponse(*registered))
}

func (h *UserHandlerImpl) Login(c *gin.Context) {
	var req dto.UserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}
	ctx := c.Request.Context()
	_, refresh, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	Cookie(c, refresh, h.cookieTTL)
	c.JSON(200, gin.H{})
}

func (h *UserHandlerImpl) GetUserById(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	user, err := h.svc.GetUserById(ctx, id)
	if err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(200, toUserResponse(*user))
}

// UpdateUser - обновляет имя/email только у самого себя: id берётся не из
// тела запроса и не из URL (иначе можно было бы редактировать чужой профиль),
// а из контекста аутентифицированного запроса, куда его кладёт Auth-миддлвар
func (h *UserHandlerImpl) UpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, response.Error{Message: err.Error(), Code: "FAILED_TO_BIND"})
		return
	}
	ctx := c.Request.Context()
	userId := c.GetString("user_id")
	user := &domain.User{Id: userId}
	if err := h.svc.UpdateUser(ctx, user, req.Name, req.Email); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(200, gin.H{})
}

func (h *UserHandlerImpl) Refresh(c *gin.Context) {
	ctx := c.Request.Context()
	token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, response.Error{Message: err.Error(), Code: "MISSING_REFRESH_TOKEN"})
		return
	}
	_, refresh, err := h.svc.Refresh(ctx, token)
	if err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	Cookie(c, refresh, h.cookieTTL)
	c.JSON(200, gin.H{})
}

func (h *UserHandlerImpl) GetTokenList(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.Param("userId")
	tokens, err := h.svc.GetTokensList(ctx, userId)
	if err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(200, gin.H{"token": tokens})
}

func (h *UserHandlerImpl) RevokeToken(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	if err := h.svc.RevokeToken(ctx, id); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(200, gin.H{})
}

func (h *UserHandlerImpl) RevokeAllTokens(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.Param("userId")
	if err := h.svc.RevokeAllTokens(ctx, userId); err != nil {
		status, body := mapServiceError(err)
		c.JSON(status, body)
		return
	}
	c.JSON(200, gin.H{})
}
