package controller

import (
	"net/http"
	"strconv"

	"prom-lens-backend/internal/apps/authn/request"
	"prom-lens-backend/internal/apps/authn/service"
	pg "prom-lens-backend/internal/core/pagination"
	"prom-lens-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	user service.UserManager
}

func NewUserController(user service.UserManager) *UserController {
	return &UserController{user: user}
}

func (c *UserController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
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
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *UserController) Create(ctx *gin.Context) {
	var req request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	req.Creator = ctx.MustGet("username").(string)
	if err := c.user.Create(ctx, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *UserController) Update(ctx *gin.Context) {
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}
	var req request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.user.Update(ctx, userID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *UserController) Delete(ctx *gin.Context) {
	userID, ok := ginhelper.RequireUintParam(ctx, "userID")
	if !ok {
		return
	}
	operatorID := uint(ctx.MustGet("uid").(int))
	if err := c.user.Delete(ctx, operatorID, userID); err != nil {
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.user.ResetPassword(ctx, userID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
