package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/service"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
)

// StatefulSetController 定义控制器结构体
type StatefulSetController struct {
	statefulSet service.StatefulSetService // 使用服务接口
}

// NewStatefulSetController 创建新的 StatefulSetController 实例
func NewStatefulSetController(statefulSet service.StatefulSetService) *StatefulSetController {
	return &StatefulSetController{statefulSet: statefulSet}
}

func (c *StatefulSetController) List(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
	)
	data, err := c.statefulSet.List(ctx, id, ns)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StatefulSetController) Get(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	data, err := c.statefulSet.Get(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StatefulSetController) GetDetail(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	data, err := c.statefulSet.GetDetail(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StatefulSetController) Create(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		body  = &appsv1.StatefulSet{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
	}
	// 创建
	data, err := c.statefulSet.Create(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StatefulSetController) Update(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
		body  = &appsv1.StatefulSet{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
	}
	// 更新
	body.Name = name
	body.Namespace = ns
	data, err := c.statefulSet.Update(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StatefulSetController) Delete(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	err := c.statefulSet.Delete(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *StatefulSetController) Scale(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
		scale Scale
	)
	if err := ctx.ShouldBindJSON(&scale); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新副本
	data, err := c.statefulSet.Scale(ctx, id, ns, name, scale.Replicas)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StatefulSetController) Restart(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	// 重启
	data, err := c.statefulSet.Restart(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
