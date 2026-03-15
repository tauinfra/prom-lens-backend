package controller

import (
	"net/http"
	"valyria-backend/internal/apps/dragon/service"

	"github.com/gin-gonic/gin"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type TektonTaskController struct {
	manager service.TektonTaskManager
}

func NewTektonTaskController(manager service.TektonTaskManager) *TektonTaskController {
	return &TektonTaskController{manager: manager}
}

func (c *TektonTaskController) List(ctx *gin.Context) {
	data, err := c.manager.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *TektonTaskController) Get(ctx *gin.Context) {
	name := ctx.Param("name")
	data, err := c.manager.Get(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *TektonTaskController) Create(ctx *gin.Context) {
	body := &tektonv1.Task{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	_, err := c.manager.Create(ctx, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *TektonTaskController) Update(ctx *gin.Context) {
	name := ctx.Param("name")
	body := &tektonv1.Task{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	_, err := c.manager.Update(ctx, name, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *TektonTaskController) Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	if err := c.manager.Delete(ctx, name); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
