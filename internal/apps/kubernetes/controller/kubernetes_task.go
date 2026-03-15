package controller

import (
	"net/http"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// TaskController 定义控制器结构体
type TaskController struct {
	task service.TaskManager // 使用服务接口
}

// NewTaskController 创建新的 TaskController 实例
func NewTaskController(task service.TaskManager) *TaskController {
	return &TaskController{task: task}
}

func (c *TaskController) List(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	ns := ctx.Param("namespace")
	data, err := c.task.List(ctx, id, ns)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *TaskController) Get(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	ns := ctx.Param("namespace")
	name := ctx.Param("name")
	data, err := c.task.Get(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *TaskController) Create(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	ns := ctx.Param("namespace")
	body := &v1.Task{}
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 创建
	_, err := c.task.Create(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *TaskController) Update(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	ns := ctx.Param("namespace")
	name := ctx.Param("name")
	body := &v1.Task{}
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新
	body.Name = name
	body.Namespace = ns
	_, err := c.task.Update(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *TaskController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	ns := ctx.Param("namespace")
	name := ctx.Param("name")
	err := c.task.Delete(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *TaskController) DeleteBatch(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	ns := ctx.Param("namespace")
	deleteBatchByNames(ctx, func(name string) error {
		return c.task.Delete(ctx, id, ns, name)
	})
}
