package controller

import (
	"strings"
	"valyria-backend/internal/apps/kubernetes/request"

	"github.com/gin-gonic/gin"
)

func deleteBatchByNames(ctx *gin.Context, deleteFn func(name string) error) {
	var req request.BatchDeleteWorkloadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	seen := make(map[string]struct{}, len(req.Names))
	for _, name := range req.Names {
		name = strings.TrimSpace(name)
		if name == "" {
			ctx.JSON(200, gin.H{"success": false, "code": -1, "msg": "name is required"})
			return
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		if err := deleteFn(name); err != nil {
			ctx.JSON(200, gin.H{"success": false, "code": -1, "msg": err.Error()})
			return
		}
	}
	ctx.JSON(200, gin.H{"success": true, "code": 200})
}

func deleteBatchByIDs(ctx *gin.Context, deleteFn func(id uint) error) {
	var req request.BatchDeleteIDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	seen := make(map[uint]struct{}, len(req.IDs))
	for _, id := range req.IDs {
		if id == 0 {
			ctx.JSON(200, gin.H{"success": false, "code": -1, "msg": "id is required"})
			return
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if err := deleteFn(id); err != nil {
			ctx.JSON(200, gin.H{"success": false, "code": -1, "msg": err.Error()})
			return
		}
	}
	ctx.JSON(200, gin.H{"success": true, "code": 200})
}
