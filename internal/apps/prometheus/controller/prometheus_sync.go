package controller

import (
	"net/http"
	"prom-lens-backend/internal/apps/prometheus/service"

	"github.com/gin-gonic/gin"
)

type SyncController struct {
	sync service.SyncManager
}

func NewSyncController(sync service.SyncManager) *SyncController {
	return &SyncController{sync: sync}
}

// ImportRules 从规则 ConfigMap 全量导入（仅新增，区分 alert 与 record）。
func (c *SyncController) ImportRules(ctx *gin.Context) {
	result, err := c.sync.ImportRulesFromConfigMap(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	success := len(result.Errors) == 0
	code := 200
	msg := "import completed"
	if !success {
		msg = "import completed with errors"
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": success,
		"code":    code,
		"msg":     msg,
		"data":    result,
	})
}
