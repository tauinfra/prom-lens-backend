package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"valyria-backend/internal/core/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Claims JWT 声明结构体
type Claims struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Su       bool   `json:"su"`
	jwt.RegisteredClaims
}

// JWTManager JWT 管理器
type JWTManager struct {
	secret              []byte
	accessTokenExpires  time.Duration
	refreshTokenExpires time.Duration
	whitelist           []string
	issuer              string
	audience            string
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{
		secret:              []byte(cfg.App.JWTSecret),
		accessTokenExpires:  cfg.Auth.AccessTokenExpires,
		refreshTokenExpires: cfg.Auth.RefreshTokenExpires,
		whitelist:           cfg.Auth.Whitelist,
		issuer:              cfg.Auth.Issuer,
		audience:            cfg.Auth.Audience,
	}
}

// TokenResponse Token 响应结构
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
}

// GenToken 生成 Token
func (j *JWTManager) GenToken(uid int, username string, isSuperuser bool) (*TokenResponse, error) {
	now := time.Now()
	accessTokenExpires := now.Add(j.accessTokenExpires)
	refreshTokenExpires := now.Add(j.refreshTokenExpires)

	// Access Token
	accessClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(accessTokenExpires),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Subject:   "access_token",
	}
	if j.issuer != "" {
		accessClaims.Issuer = j.issuer
	}
	if j.audience != "" {
		accessClaims.Audience = jwt.ClaimStrings{j.audience}
	}
	accessTokenClaims := Claims{
		UserID:           uid,
		Username:         username,
		Su:               isSuperuser,
		RegisteredClaims: accessClaims,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(j.secret)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshTokenExpires),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Subject:   "refresh_token",
	}
	if j.issuer != "" {
		refreshClaims.Issuer = j.issuer
	}
	if j.audience != "" {
		refreshClaims.Audience = jwt.ClaimStrings{j.audience}
	}
	refreshTokenClaims := Claims{
		UserID:           uid,
		Username:         username,
		Su:               isSuperuser,
		RegisteredClaims: refreshClaims,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(j.secret)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    accessTokenExpires.Format("2006/01/02 15:04:05"),
	}, nil
}

// ParseToken 解析 Token
func (j *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid or expired token")
}

// RefreshToken 刷新 Token
func (j *JWTManager) RefreshToken(refreshToken string) (*TokenResponse, error) {
	claims, err := j.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// 验证这是 refresh token
	if claims.Subject != "refresh_token" {
		return nil, errors.New("invalid token type")
	}

	return j.GenToken(claims.UserID, claims.Username, claims.Su)
}

// extractToken 从请求中提取 Token
func (j *JWTManager) extractToken(c *gin.Context) string {
	// 1. 尝试从 Authorization Header 获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return parts[1]
		}
	}

	// 2. 尝试从 WebSocket Protocol 获取
	if token := c.GetHeader("Sec-WebSocket-Protocol"); token != "" {
		return token
	}

	return ""
}

// Cors 跨域
func (j *JWTManager) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")                                                       // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization") // 你想放行的 header 也可以在后面自行添加
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE, UPDATE")         //
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// RequireAuth 认证中间件
func (j *JWTManager) RequireAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 检查白名单
		if HasWhiteList(c.Request.URL, j.whitelist) {
			c.Next()
			return
		}
		// 提取 Token
		token := j.extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    401,
				"msg":     "请求未授权，缺少认证信息",
			})
			c.Abort()
			return
		}
		// 解析 Token
		claims, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    401,
				"msg":     "Token 无效或已过期，请重新登录",
			})
			c.Abort()
			return
		}
		// 校验 issuer/audience（配置为空时不校验）
		if j.issuer != "" && claims.Issuer != j.issuer {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    401,
				"msg":     "Token 签发方无效，请重新登录",
			})
			c.Abort()
			return
		}
		if j.audience != "" && !claims.VerifyAudience(j.audience, true) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    401,
				"msg":     "Token 受众无效，请重新登录",
			})
			c.Abort()
			return
		}
		// 仅允许 access token 通过认证中间件
		if claims.Subject != "access_token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    401,
				"msg":     "Token 类型无效，请重新登录",
			})
			c.Abort()
			return
		}
		// 设置用户信息到上下文
		c.Set("uid", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("isSuperuser", claims.Su) // 超级管理员
		c.Next()
	}
}
