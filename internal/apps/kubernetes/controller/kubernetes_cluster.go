package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/request"
	"valyria-backend/internal/apps/kubernetes/service"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

// ClusterController 定义控制器结构体
type ClusterController struct {
	cluster service.ClusterManager // 使用服务接口
}

func NewClusterController(cluster service.ClusterManager) *ClusterController {
	return &ClusterController{cluster: cluster}
}

func (c *ClusterController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"), // 搜索关键字
	}

	data, pagination, err := c.cluster.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *ClusterController) Get(ctx *gin.Context) {
	clusterIDU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	clusterID := clusterIDU
	data, err := c.cluster.Get(ctx, clusterID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ClusterController) Create(ctx *gin.Context) {
	var req request.CreateClusterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 从 context 获取用户名并传递给 service
	var username string
	if usernameValue, exists := ctx.Get("username"); exists {
		if usernameStr, ok := usernameValue.(string); ok && usernameStr != "" {
			username = usernameStr
		}
	}
	if err := c.cluster.Create(ctx, &req, username); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ClusterController) Update(ctx *gin.Context) {
	var req request.UpdateClusterRequest
	clusterIDU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	clusterID := clusterIDU

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.cluster.Update(ctx, clusterID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ClusterController) UpdateToken(ctx *gin.Context) {
	var req request.UpdateClusterTokenRequest
	clusterIDU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	clusterID := clusterIDU
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.cluster.UpdateToken(ctx, clusterID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ClusterController) Delete(ctx *gin.Context) {
	clusterIDU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	clusterID := clusterIDU

	if err := c.cluster.Delete(ctx, clusterID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
