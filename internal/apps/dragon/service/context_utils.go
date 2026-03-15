package service

import (
	"context"

	"github.com/gin-gonic/gin"
)

func IsSuperuser(ctx context.Context) bool {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return false
	}
	val, exists := ginCtx.Get("isSuperuser")
	if !exists {
		return false
	}
	isSuper, ok := val.(bool)
	return ok && isSuper
}

func GetUID(ctx context.Context) uint {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return 0
	}
	val, exists := ginCtx.Get("uid")
	if !exists {
		return 0
	}
	if id, ok := val.(int); ok && id > 0 {
		return uint(id)
	}
	return 0
}

func GetUsername(ctx context.Context) string {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return ""
	}
	return ginCtx.GetString("username")
}
