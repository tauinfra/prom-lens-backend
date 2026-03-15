package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/authn/request"
	"valyria-backend/internal/apps/authn/service"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

// UserController 定义控制器结构体
type UserController struct {
	user service.UserService // 使用服务接口
}

func NewUserController(user service.UserService) *UserController {
	return &UserController{
		user: user,
	}
}

func (c *UserController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
	}

	data, pagination, err := c.user.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *UserController) Get(ctx *gin.Context) {
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}
	data, err := c.user.Get(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *UserController) Create(ctx *gin.Context) {
	var req request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	username := ctx.MustGet("username").(string)
	req.Creator = username
	if err := c.user.Create(ctx, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *UserController) Update(ctx *gin.Context) {
	var req request.UpdateUserRequest
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.user.Update(ctx, userID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *UserController) Delete(ctx *gin.Context) {
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}
	if err := c.user.Delete(ctx, userID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *UserController) ResetPassword(ctx *gin.Context) {
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}
	var req request.ResetPasswordRequest
	// 绑定结构体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 重置密码
	if err := c.user.ResetPassword(ctx, userID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *UserController) GetUserRoles(ctx *gin.Context) {
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}
	data, err := c.user.GetUserRoles(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *UserController) UpdateRoleRoles(ctx *gin.Context) {
	var req request.UpdateUserRolesRequest
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.user.UpdateUserRoles(ctx, userID, req.Roles); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *UserController) GetUserMenus(ctx *gin.Context) {
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}
	data, err := c.user.GetUserMenus(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *UserController) UpdateUserMenus(ctx *gin.Context) {
	var req request.UpdateUserMenusRequest
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.user.UpdateUserMenus(ctx, userID, req.Menus); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
