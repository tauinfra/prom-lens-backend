package controller

import (
	"net/http"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
	networkingv1 "k8s.io/api/networking/v1"
)

type IngressClassController struct {
	ingressClass service.IngressClassManager
}

func NewIngressClassController(ingressClass service.IngressClassManager) *IngressClassController {
	return &IngressClassController{ingressClass: ingressClass}
}

func (c *IngressClassController) List(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	data, err := c.ingressClass.List(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *IngressClassController) Get(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	data, err := c.ingressClass.Get(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *IngressClassController) Create(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	body := &networkingv1.IngressClass{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	_, err := c.ingressClass.Create(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *IngressClassController) Update(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	body := &networkingv1.IngressClass{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	body.Name = name
	_, err := c.ingressClass.Update(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *IngressClassController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	if err := c.ingressClass.Delete(ctx, id, name); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *IngressClassController) DeleteBatch(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	deleteBatchByNames(ctx, func(name string) error {
		return c.ingressClass.Delete(ctx, id, name)
	})
}
