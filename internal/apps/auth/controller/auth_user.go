package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/auth/model"
	"valyria-backend/internal/apps/auth/service"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/core/middleware"
	pg "valyria-backend/internal/core/pagination"

	"github.com/gin-gonic/gin"
)

// UserController 定义控制器结构体
type UserController struct {
	jwtManager *middleware.JWTManager
	user       service.UserService // 使用服务接口
}

func NewUserController(user service.UserService, jwtManager *middleware.JWTManager) *UserController {
	return &UserController{
		user:       user,
		jwtManager: jwtManager,
	}
}

// login 登录
type login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// changePassword 更改密码
type changePassword struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"` // 可以设置新密码的最小长度、复杂度等策略
}

// resetPassword 重置密码
type resetPassword struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"` // 可以设置新密码的最小长度、复杂度等策略
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
	}

	data, pagination, err := c.user.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *UserController) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	data, err := c.user.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *UserController) Create(ctx *gin.Context) {
	var data model.User
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.user.Create(ctx, &data); err != nil {
		logger.Errorf("User '%v' creation failed, error:: %v", data.Username, err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	logger.Infof("User '%v' has been created successfully.", data.Username)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *UserController) Update(ctx *gin.Context) {
	var data model.User
	var id, _ = strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	fmt.Println(data.RoleIDs)
	if err := c.user.Update(ctx, id, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	logger.Infof("User '%v' has been updated successfully.", data.Username)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *UserController) Delete(ctx *gin.Context) {
	var id, _ = strconv.Atoi(ctx.Param("id"))

	if err := c.user.Delete(ctx, id); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	logger.Infof("UserID '%d' has been deleted successfully.", id)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *UserController) Login(ctx *gin.Context) {
	var login login
	// 绑定结构体
	if err := ctx.ShouldBindJSON(&login); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 登录
	if user, err := c.user.Login(ctx, login.Username, login.Password); err != nil {
		logger.Errorf("User '%v' failed to log in, error: %v", login.Username, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": http.StatusUnauthorized, "msg": err.Error()})
		return
	} else {
		res, _ := c.jwtManager.GenToken(user.ID, user.Username) // 生成 accessToken 和 refreshToken 信息
		roles := []string{}
		for _, v := range user.Roles {
			roles = append(roles, v.Name)
		}
		ctx.JSON(http.StatusOK, gin.H{ // 用户信息(根据前端框架生成固定格式)
			"success": true,
			"data": gin.H{
				"uid":          user.ID,
				"username":     user.Username,
				"accessToken":  res.AccessToken,
				"refreshToken": res.RefreshToken,
				"expires":      res.ExpiresIn,
				"roles":        roles,
			},
		})
		ctx.Next()
	}
}
