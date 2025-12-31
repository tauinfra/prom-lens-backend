package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/service"

	"github.com/gin-gonic/gin"
)

// EventController 定义控制器结构体
type EventController struct {
	event service.EventService // 使用服务接口
}

// NewEventController 创建新的 EventController 实例
func NewEventController(event service.EventService) *EventController {
	return &EventController{
		event: event,
	}
}

func (c *EventController) List(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		kind  = ctx.Query("kind")
		name  = ctx.Query("name")
	)
	data, err := c.event.List(ctx, id, ns, kind, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
