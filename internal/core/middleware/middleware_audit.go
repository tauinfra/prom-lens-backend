package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	model2 "valyria-backend/internal/apps/audit/model"
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

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// MaskSensitiveFields 递归替换敏感字段
func (a *AuditManager) MaskSensitiveFields(data map[string]interface{}) {
	var sensitiveKeys = []string{"password", "token", "secret", "key"}
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

// ParseAndMaskJSON 转换并掩码
func (a *AuditManager) ParseAndMaskJSON(raw []byte) datatypes.JSON {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(raw, &jsonData); err != nil {
		// 解析失败，直接返回原始内容
		return datatypes.JSON(raw)
	}
	a.MaskSensitiveFields(jsonData)
	masked, _ := json.Marshal(jsonData)
	return datatypes.JSON(masked)
}

// AuditLogin 登录审计
func (a *AuditManager) AuditLogin(tx *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" || !strings.Contains(c.FullPath(), "login") {
			c.Next()
			return
		}

		params, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(params))

		var loginR login
		if err := json.Unmarshal(params, &loginR); err != nil {
			logger.Errorf("login body unmarshal failed: %v", err)
			c.Next()
			return
		}

		ua := useragent.Parse(c.GetHeader("User-Agent"))
		authLog := model2.AuthLog{
			Username:   loginR.Username,
			IPAddress:  c.ClientIP(),
			System:     ua.OS,
			Agent:      ua.Name + " " + ua.Version,
			StatusCode: c.Writer.Status(),
		}

		if err := tx.Create(&authLog).Error; err != nil {
			logger.Errorf("write auth log failed: %v", err)
		}

		c.Next()
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
		params, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(params))

		// 包装 Response
		res := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = res

		// 执行业务逻辑
		c.Next()

		// 提取用户信息
		username, exists := c.Get("username")
		if !exists {
			// 用户未登录，直接跳过记录
			c.Next()
			return
		}
		// 如果用户信息存在，则记录操作日志
		ua := useragent.Parse(c.GetHeader("User-Agent"))
		oplog := model2.AuditLog{
			Username:   username.(string),
			IPAddress:  c.ClientIP(),
			Method:     c.Request.Method,
			UrlPath:    c.Request.URL.String(),
			Agent:      ua.Name + " " + ua.Version,
			StatusCode: c.Writer.Status(),
			Params:     a.ParseAndMaskJSON(params),
			Response:   a.ParseAndMaskJSON(res.body.Bytes()),
		}
		if err := tx.Create(&oplog).Error; err != nil {
			logger.Errorf("write operation log failed: %v", err)
		}
	}
}
