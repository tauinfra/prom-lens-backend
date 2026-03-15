package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"
	"valyria-backend/internal/apps/audit/model"
	"valyria-backend/internal/core/logger"

	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AuditManager 审计管理器
type AuditManager struct {
	db *gorm.DB
}

// NewAuditManager 创建审计管理器
func NewAuditManager(db *gorm.DB) *AuditManager {
	return &AuditManager{
		db: db,
	}
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResult struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
}

type responseWriter struct {
	gin.ResponseWriter
	body         *bytes.Buffer
	maxBodySize  int
	bodyTruncate bool
}

func (r *responseWriter) Write(b []byte) (int, error) {
	if r.maxBodySize > 0 && r.body.Len() < r.maxBodySize {
		remaining := r.maxBodySize - r.body.Len()
		if len(b) > remaining {
			r.body.Write(b[:remaining])
			r.bodyTruncate = true
		} else {
			r.body.Write(b)
		}
	} else if r.maxBodySize > 0 {
		r.bodyTruncate = true
	}
	return r.ResponseWriter.Write(b)
}

const maxAuditBodySize = 1 << 20 // 1MB
const auditWriteTimeout = 2 * time.Second
const maxAuditStringSize = 255

func readBodyWithLimit(r io.Reader, limit int) ([]byte, bool) {
	if limit <= 0 {
		data, _ := io.ReadAll(r)
		return data, false
	}
	data, _ := io.ReadAll(io.LimitReader(r, int64(limit+1)))
	if len(data) > limit {
		return data[:limit], true
	}
	return data, false
}

func truncateString(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max]
}

func parseLoginResult(contentType string, body []byte, statusCode int) (success bool, code int, msg string) {
	if statusCode >= 200 && statusCode < 400 {
		success = true
	}
	if !isJSONContentType(contentType) || len(body) == 0 {
		return success, 0, ""
	}
	var res loginResult
	if err := json.Unmarshal(body, &res); err != nil {
		return success, 0, ""
	}
	if res.Success {
		return true, res.Code, ""
	}
	return false, res.Code, res.Msg
}

func writeAuditLog(tx *gorm.DB, record interface{}, logLabel string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), auditWriteTimeout)
		defer cancel()
		if err := tx.WithContext(ctx).Create(record).Error; err != nil {
			logger.Errorf("write %s log failed: %v", logLabel, err)
		}
	}()
}

// MaskSensitiveFields 递归替换敏感字段
func (a *AuditManager) MaskSensitiveFields(data map[string]interface{}) {
	var sensitiveKeys = []string{
		"password", "token", "secret", "key",
		"authorization", "apikey", "api_key", "apisecret", "api_secret",
		"access_token", "refresh_token", "private_key", "credential",
	}
	for key, value := range data {
		lowerKey := strings.ToLower(key)
		for _, sKey := range sensitiveKeys {
			if strings.Contains(lowerKey, sKey) {
				data[key] = "********"
			}
		}
		switch v := value.(type) {
		case map[string]interface{}:
			a.MaskSensitiveFields(v)
		case []interface{}:
			for _, item := range v {
				if nested, ok := item.(map[string]interface{}); ok {
					a.MaskSensitiveFields(nested)
				}
			}
		}
	}
}

func isJSONContentType(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "application/*+json")
}

func formatUserAgent(ua useragent.UserAgent) (system, agent string) {
	system = strings.TrimSpace(ua.OS)
	agent = strings.TrimSpace(ua.Name + " " + ua.Version)
	if system == "" {
		system = "unknown"
	}
	if agent == "" {
		agent = "unknown"
	}
	return system, agent
}

// ParseAndMaskJSON 转换并掩码
func (a *AuditManager) ParseAndMaskJSON(raw []byte) datatypes.JSON {
	if len(raw) == 0 {
		return datatypes.JSON(raw)
	}
	var jsonData map[string]interface{}
	if err := json.Unmarshal(raw, &jsonData); err != nil {
		// 解析失败，返回空内容避免泄露
		return datatypes.JSON([]byte("{}"))
	}
	a.MaskSensitiveFields(jsonData)
	masked, _ := json.Marshal(jsonData)
	return datatypes.JSON(masked)
}

// AuditLogin 登录审计
func (a *AuditManager) AuditLogin(tx *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" || c.FullPath() != "/api/v1/authn/login" {
			c.Next()
			return
		}

		params, _ := readBodyWithLimit(c.Request.Body, maxAuditBodySize)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(params))

		// 包装 Response，确保拿到真实状态码
		res := &responseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
			maxBodySize:    maxAuditBodySize,
		}
		c.Writer = res

		var loginR login
		if err := json.Unmarshal(params, &loginR); err != nil {
			logger.Errorf("login body unmarshal failed: %v", err)
			c.Next()
			return
		}

		c.Next()

		ua := useragent.Parse(c.GetHeader("User-Agent"))
		system, agent := formatUserAgent(ua)
		success, code, msg := parseLoginResult(c.Writer.Header().Get("Content-Type"), res.body.Bytes(), c.Writer.Status())
		authLog := model.AuthLog{
			Username:   truncateString(loginR.Username, maxAuditStringSize),
			IPAddress:  c.ClientIP(),
			System:     truncateString(system, maxAuditStringSize),
			Agent:      truncateString(agent, maxAuditStringSize),
			StatusCode: c.Writer.Status(),
			Success:    success,
			ErrorCode:  code,
			ErrorMsg:   truncateString(msg, maxAuditStringSize),
		}
		writeAuditLog(tx, &authLog, "authn")
	}
}

// AuditOperation 操作审计
func (a *AuditManager) AuditOperation(tx *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过 GET
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 读取 Request
		params, _ := readBodyWithLimit(c.Request.Body, maxAuditBodySize)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(params))

		// 包装 Response
		res := &responseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
			maxBodySize:    maxAuditBodySize,
		}
		c.Writer = res

		// 执行业务逻辑
		c.Next()

		// 提取用户信息
		username, exists := c.Get("username")
		if !exists {
			// 用户未登录，直接跳过记录
			return
		}
		// 如果用户信息存在，则记录操作日志
		usernameStr, ok := username.(string)
		if !ok || usernameStr == "" {
			return
		}
		ua := useragent.Parse(c.GetHeader("User-Agent"))
		_, agent := formatUserAgent(ua)
		var paramsJSON, responseJSON datatypes.JSON
		if isJSONContentType(c.GetHeader("Content-Type")) {
			paramsJSON = a.ParseAndMaskJSON(params)
		} else {
			paramsJSON = datatypes.JSON([]byte("{}"))
		}
		if isJSONContentType(c.Writer.Header().Get("Content-Type")) {
			responseJSON = a.ParseAndMaskJSON(res.body.Bytes())
		} else {
			responseJSON = datatypes.JSON([]byte("{}"))
		}

		oplog := model.AuditLog{
			Username:   truncateString(usernameStr, maxAuditStringSize),
			IPAddress:  truncateString(c.ClientIP(), maxAuditStringSize),
			Method:     c.Request.Method,
			UrlPath:    truncateString(c.Request.URL.String(), maxAuditStringSize),
			Agent:      truncateString(agent, maxAuditStringSize),
			StatusCode: c.Writer.Status(),
			Success:    c.Writer.Status() >= 200 && c.Writer.Status() < 400,
			Params:     paramsJSON,
			Response:   responseJSON,
		}
		writeAuditLog(tx, &oplog, "operation")
	}
}
