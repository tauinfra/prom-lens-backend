package ginhelper

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseUintParam 从路径参数解析无符号正整数，用于 ID 类参数校验。
// 仅当参数存在且为大于 0 的整数时返回 (value, true)，否则返回 (0, false)。
func ParseUintParam(ctx *gin.Context, key string) (uint, bool) {
	s := ctx.Param(key)
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, false
	}
	return uint(n), true
}

// RequireUintParam 解析路径参数为 uint；若无效则写入 400 并 Abort，返回 (0, false)。
// 用法：id, ok := ginhelper.RequireUintParam(ctx, "projectID"); if !ok { return }
func RequireUintParam(ctx *gin.Context, key string) (uint, bool) {
	id, ok := ParseUintParam(ctx, key)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "code": 400, "msg": "invalid or missing " + key})
		ctx.Abort()
		return 0, false
	}
	return id, true
}
