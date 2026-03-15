package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"valyria-backend/internal/apps/authn/dto"
	"valyria-backend/internal/apps/authn/request"
	"valyria-backend/internal/apps/authn/service"
	"valyria-backend/internal/core/config"
	"valyria-backend/internal/core/database"
	"valyria-backend/internal/core/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// AuthController 定义控制器结构体
type AuthController struct {
	jwtManager *middleware.JWTManager
	auth       service.AuthManager // 使用服务接口
}

func NewAuthController(auth service.AuthManager, jwtManager *middleware.JWTManager) *AuthController {
	return &AuthController{
		auth:       auth,
		jwtManager: jwtManager,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	var roles []string
	var permissions []string
	// 绑定结构体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 登录保护：失败次数限制（按 IP 与用户名）
	maxRetries, lockWindow := getLoginProtectionConfig()
	if maxRetries > 0 && lockWindow > 0 && database.Redis != nil {
		if locked, wait := isLoginLocked(ctx, req.Username, maxRetries, lockWindow); locked {
			retryAfter := int(wait.Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"code":    http.StatusTooManyRequests,
				"msg":     fmt.Sprintf("登录失败次数过多，账号已被锁定，请联系运维人员处理"),
			})
			return
		}
	}
	// 登录
	user, err := c.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		if maxRetries > 0 && lockWindow > 0 && database.Redis != nil {
			recordLoginFailure(ctx, req.Username, maxRetries, lockWindow)
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": http.StatusUnauthorized, "msg": err.Error()})
		return
	}
	if maxRetries > 0 && lockWindow > 0 && database.Redis != nil {
		clearLoginFailures(ctx, req.Username)
	}
	// 登录成功：更新最后登录时间后再生成 JWT
	_ = c.auth.UpdateLastLogin(ctx, uint(user.ID))
	isSuperuser := user.IsSuperuser != nil && *user.IsSuperuser
	res, _ := c.jwtManager.GenToken(user.ID, user.Username, isSuperuser) // 生成 accessToken 和 refreshToken 信息
	if isSuperuser {
		roles = append(roles, "admin")
		permissions = append(permissions, "*:*:*")
	} else {
		roles, permissions, err = c.auth.GetRolePermissions(ctx, uint(user.ID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{ // 用户信息(根据前端框架生成固定格式)
		"success": true,
		"data": dto.LoginResponseDTO{
			UID:          user.ID,
			Username:     user.Username,
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			Expires:      res.ExpiresIn,
			Roles:        roles,
			Permissions:  permissions,
		},
	})
	ctx.Next()
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "缺少 Authorization Header",
		})
		return
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "Authorization 格式错误",
		})
		return
	}

	refreshToken := strings.TrimPrefix(authHeader, prefix)
	if refreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"msg":     "refresh token 不能为空",
		})
		return
	}
	// 刷新 Token
	res, err := c.jwtManager.RefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": err.Error()})
		return
	}
	// 生成 accessToken 和 refreshToken
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"data": dto.LoginResponseDTO{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			Expires:      res.ExpiresIn,
		},
	})
	ctx.Next()
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req request.ChangePasswordRequest
	// 绑定结构体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	uid := ctx.MustGet("uid") // 假设用户身份信息存储在上下文中
	// 修改密码
	if err := c.auth.ChangePassword(ctx, uint(uid.(int)), &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *AuthController) GetUserRoutes(ctx *gin.Context) {
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
	isSuperuserVal, ok := ctx.Get("isSuperuser")
	isSuperuser := false
	if ok {
		if v, ok := isSuperuserVal.(bool); ok {
			isSuperuser = v
		}
	}
	data, err := c.auth.GetUserRoutes(ctx, uint(uid), isSuperuser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

// GetProfile 当前用户 profile GET /authn/me/profile
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
	isSuperuserVal, _ := ctx.Get("isSuperuser")
	isSuperuser := false
	if v, ok := isSuperuserVal.(bool); ok {
		isSuperuser = v
	}
	var roles, permissions []string
	if isSuperuser {
		roles = []string{"Admin"}
		permissions = []string{"*:*:*"}
	} else {
		roles, permissions, err = c.auth.GetRolePermissions(ctx, uint(uid))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "code": -1, "msg": err.Error()})
			return
		}
		if roles == nil {
			roles = []string{}
		}
		if permissions == nil {
			permissions = []string{}
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    200,
		"data": dto.ProfileDTO{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Phone:       user.Phone,
			CreatedAt:   user.CreatedAt,
			LastLogin:   user.LastLogin,
			Roles:       roles,
			Permissions: permissions,
		},
	})
}

func getLoginProtectionConfig() (int, time.Duration) {
	cfg := config.GetConfig()
	if cfg == nil {
		return 0, 0
	}
	return cfg.Auth.LoginMaxRetries, cfg.Auth.LoginLockWindow
}

func loginKeyByIP(ip string) string {
	return fmt.Sprintf("authn:login:ip:%s", ip)
}

func loginKeyByUser(username string) string {
	return fmt.Sprintf("authn:login:user:%s", username)
}

func isLoginLocked(ctx *gin.Context, username string, maxRetries int, lockWindow time.Duration) (bool, time.Duration) {
	ip := ctx.ClientIP()
	keys := []string{loginKeyByIP(ip)}
	if username != "" {
		keys = append(keys, loginKeyByUser(username))
	}
	for _, key := range keys {
		count, err := database.Redis.Get(ctx, key).Int()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			// Redis 异常时跳过保护，避免误伤
			continue
		}
		if count >= maxRetries {
			ttl, err := database.Redis.TTL(ctx, key).Result()
			if err != nil {
				return true, lockWindow
			}
			return true, ttl
		}
	}
	return false, 0
}

func recordLoginFailure(ctx *gin.Context, username string, maxRetries int, lockWindow time.Duration) {
	ip := ctx.ClientIP()
	keys := []string{loginKeyByIP(ip)}
	if username != "" {
		keys = append(keys, loginKeyByUser(username))
	}
	pipe := database.Redis.Pipeline()
	for _, key := range keys {
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, lockWindow)
	}
	_, _ = pipe.Exec(ctx)
}

func clearLoginFailures(ctx *gin.Context, username string) {
	ip := ctx.ClientIP()
	keys := []string{loginKeyByIP(ip)}
	if username != "" {
		keys = append(keys, loginKeyByUser(username))
	}
	_ = database.Redis.Del(ctx, keys...).Err()
}
