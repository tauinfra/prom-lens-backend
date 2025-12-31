package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/model"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/core/logger"
	pg "valyria-backend/internal/core/pagination"

	"github.com/gin-gonic/gin"
)

// ClusterController 定义控制器结构体
type ClusterController struct {
	cluster service.ClusterService // 使用服务接口
}

func NewClusterController(cluster service.ClusterService) *ClusterController {
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
	id, _ := strconv.Atoi(ctx.Param("id"))
	data, err := c.cluster.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ClusterController) Create(ctx *gin.Context) {
	var data model.Cluster
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.cluster.Create(ctx, &data); err != nil {
		logger.Errorf("Kubernetes Cluster '%v' creation failed, error:: %v", data.Name, err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	logger.Infof("Kubernetes Cluster '%v' has been created successfully.", data.Name)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ClusterController) Update(ctx *gin.Context) {
	var data model.Cluster
	var id, _ = strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.cluster.Update(ctx, id, &data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	logger.Infof("Kubernetes Cluster '%v' has been updated successfully.", data.Name)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ClusterController) Delete(ctx *gin.Context) {
	var id, _ = strconv.Atoi(ctx.Param("id"))

	if err := c.cluster.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	logger.Infof("Kubernetes ClusterID '%d' has been deleted successfully.", id)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
