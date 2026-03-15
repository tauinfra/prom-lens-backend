package controller

import (
	"net/http"
	"valyria-backend/internal/apps/authn/request"
	"valyria-backend/internal/apps/authn/service"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

// MenuController 定义控制器结构体
type MenuController struct {
	Menu service.MenuManager
}

func NewMenuController(menu service.MenuManager) *MenuController {
	return &MenuController{Menu: menu}
}

func (c *MenuController) List(ctx *gin.Context) {
	data, err := c.Menu.ListTree(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *MenuController) Get(ctx *gin.Context) {
	menuID, ok := ginhelper.RequireUintParam(ctx, "menuID")
	if !ok {
		return
	}
	data, err := c.Menu.Get(ctx, menuID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *MenuController) Create(ctx *gin.Context) {
	var req request.CreateMenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	username := ctx.MustGet("username").(string)
	req.Creator = username
	if err := c.Menu.Create(ctx, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *MenuController) Update(ctx *gin.Context) {
	var req request.UpdateMenuRequest
	menuID, ok := ginhelper.RequireUintParam(ctx, "menuID")
	if !ok {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.Menu.Update(ctx, menuID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *MenuController) Delete(ctx *gin.Context) {
	menuID, ok := ginhelper.RequireUintParam(ctx, "menuID")
	if !ok {
		return
	}

	if err := c.Menu.Delete(ctx, menuID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
