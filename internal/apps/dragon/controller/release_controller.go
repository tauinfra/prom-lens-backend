package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"valyria-backend/internal/apps/dragon/request"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/core/database"
	"valyria-backend/internal/core/logger"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/ginhelper"
	"valyria-backend/internal/pkg/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ReleaseController 定义控制器结构体
type ReleaseController struct {
	release service.ReleaseManager // 使用服务接口
}

func NewReleaseController(release service.ReleaseManager) *ReleaseController {
	return &ReleaseController{release: release}
}

func (c *ReleaseController) List(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
		Preloads:  []string{"Pipeline.Environment.Project"},
	}
	data, pagination, err := c.release.List(ctx, projectID, environmentID, pipelineID, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *ReleaseController) Get(ctx *gin.Context) {
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	releaseID, ok := ginhelper.RequireUintParam(ctx, "releaseID")
	if !ok {
		return
	}
	data, err := c.release.Get(ctx, pipelineID, releaseID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ReleaseController) Create(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	var req request.CreateReleaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	req.ProjectID = projectID
	req.EnvironmentID = environmentID
	req.PipelineID = pipelineID
	req.Creator = ctx.GetString("username")
	data, err := c.release.Create(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ReleaseController) Rollback(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	var req request.RollbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.release.Rollback(ctx, projectID, environmentID, pipelineID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ReleaseController) Delete(ctx *gin.Context) {
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	releaseID, ok := ginhelper.RequireUintParam(ctx, "releaseID")
	if !ok {
		return
	}
	if err := c.release.Delete(ctx, pipelineID, releaseID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ReleaseController) GetLogs(ctx *gin.Context) {
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	releaseID, ok := ginhelper.RequireUintParam(ctx, "releaseID")
	if !ok {
		return
	}
	data, err := c.release.Get(ctx, pipelineID, releaseID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Errorf("Failed to upgrade websocket: %v", err)
		return
	}
	defer wsConn.Close()

	conn := ws.NewWS(wsConn)
	defer conn.Close()

	lastIndex := int64(0)
	rdsKey := "release:" + data.TaskID + ":logs"
	endMsg := fmt.Sprintf("%s Release TaskID: %s completed", strings.Repeat("=", 10), data.TaskID)
	for {
		// 只取增量
		logs, err := database.Redis.LRange(ctx, rdsKey, lastIndex, -1).Result()
		if err != nil {
			logger.Errorf("Redis LRange error: %v", err)
			return
		}

		if len(logs) > 0 {
			for _, line := range logs {
				if err = conn.Write([]byte(line)); err != nil {
					logger.Infof("WebSocket disconnected: %v", err)
					return
				}
				lastIndex++
				if strings.Contains(line, endMsg) {
					logger.Infof("tekton TaskRun '%v' log reading completed.", data.TaskID)
					time.Sleep(2 * time.Second)
					return
				}
			}
		}

		// 没有新日志，阻塞等待或客户端断开
		select {
		case <-ctx.Request.Context().Done():
			logger.Info("Client disconnected or context cancelled.")
			return
		case <-time.After(1 * time.Second):
		}
	}
}
