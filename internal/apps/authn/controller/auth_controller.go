package controller

import (
	"net/http"
	"strings"
	"prom-lens-backend/internal/apps/authn/dto"
	"prom-lens-backend/internal/apps/authn/request"
	"prom-lens-backend/internal/apps/authn/service"
	"prom-lens-backend/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	jwtManager *middleware.JWTManager
	auth       service.AuthManager
}

func NewAuthController(auth service.AuthManager, jwtManager *middleware.JWTManager) *AuthController {
	return &AuthController{
		auth:       auth,
		jwtManager: jwtManager,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}

	user, err := c.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": http.StatusUnauthorized, "msg": err.Error()})
		return
	}

	_ = c.auth.UpdateLastLogin(ctx, uint(user.ID))
	isSuperuser := user.IsSuperuser != nil && *user.IsSuperuser
	res, err := c.jwtManager.GenToken(user.ID, user.Username, isSuperuser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dto.LoginResponseDTO{
			UID:          user.ID,
			Username:     user.Username,
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			Expires:      res.ExpiresIn,
		},
	})
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "缺少 Authorization Header"})
		return
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "Authorization 格式错误"})
		return
	}

	refreshToken := strings.TrimPrefix(authHeader, prefix)
	if refreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "refresh token 不能为空"})
		return
	}

	res, err := c.jwtManager.RefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"data": dto.LoginResponseDTO{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			Expires:      res.ExpiresIn,
		},
	})
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req request.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	uid := ctx.MustGet("uid")
	if err := c.auth.ChangePassword(ctx, uint(uid.(int)), &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *AuthController) GetProfile(ctx *gin.Context) {
	uidVal, ok := ctx.Get("uid")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": http.StatusUnauthorized, "msg": "请求未授权，缺少认证信息"})
		return
	}
	uid, ok := uidVal.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": http.StatusUnauthorized, "msg": "请求未授权，认证信息无效"})
		return
	}

	user, err := c.auth.GetMe(ctx, uint(uid))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"data": dto.ProfileDTO{
			ID:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
			LastLogin: user.LastLogin,
		},
	})
}
