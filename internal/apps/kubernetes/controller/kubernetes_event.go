package controller

import (
	"net/http"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

// EventController 定义控制器结构体
type EventController struct {
	event service.EventManager // 使用服务接口
}

// NewEventController 创建新的 EventController 实例
func NewEventController(event service.EventManager) *EventController {
	return &EventController{
		event: event,
	}
}

func (c *EventController) List(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	ns := ctx.Param("namespace")
	kind := ctx.Query("kind")
	name := ctx.Query("name")
	data, err := c.event.List(ctx, id, ns, kind, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
