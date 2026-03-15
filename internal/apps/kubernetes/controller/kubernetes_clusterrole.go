package controller

import (
	"net/http"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
	rbacv1 "k8s.io/api/rbac/v1"
)

type ClusterRoleController struct {
	clusterRole service.ClusterRoleManager
}

func NewClusterRoleController(clusterRole service.ClusterRoleManager) *ClusterRoleController {
	return &ClusterRoleController{clusterRole: clusterRole}
}

func (c *ClusterRoleController) List(ctx *gin.Context) {
	id, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}

	data, err := c.clusterRole.List(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ClusterRoleController) Get(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	data, err := c.clusterRole.Get(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ClusterRoleController) Create(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	body := &rbacv1.ClusterRole{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	_, err := c.clusterRole.Create(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ClusterRoleController) Update(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	body := &rbacv1.ClusterRole{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	body.Name = name
	_, err := c.clusterRole.Update(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ClusterRoleController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	if err := c.clusterRole.Delete(ctx, id, name); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ClusterRoleController) DeleteBatch(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	deleteBatchByNames(ctx, func(name string) error {
		return c.clusterRole.Delete(ctx, id, name)
	})
}
