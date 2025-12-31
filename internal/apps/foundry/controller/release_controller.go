package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/service"
	"valyria-backend/internal/core/database"
	"valyria-backend/internal/core/logger"
	pg "valyria-backend/internal/core/pagination"
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
	pipelineID, _ := strconv.Atoi(ctx.Param("pipelineID"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
		Preloads:  []string{"Pipeline.Environment.Project", "Pipeline.Application", "Pipeline.Environment.CredHarbor"},
		Filters:   map[string]interface{}{"pipeline_id": pipelineID},
	}

	data, pagination, err := c.release.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *ReleaseController) Get(ctx *gin.Context) {
	releaseID, _ := strconv.Atoi(ctx.Param("releaseID"))
	data, err := c.release.Get(ctx, releaseID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ReleaseController) Create(ctx *gin.Context) {
	var pipelineID, _ = strconv.Atoi(ctx.Param("pipelineID"))
	var data model.Release
	var username = ctx.GetString("username")
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	data.PipelineID = pipelineID // 流水线 ID
	data.DeployedBy = username   // 发布人员
	if err := c.release.Create(ctx, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ReleaseController) Update(ctx *gin.Context) {
	var data model.Release
	var pipelineID, _ = strconv.Atoi(ctx.Param("pipelineID")) // Env ID
	var releaseID, _ = strconv.Atoi(ctx.Param("releaseID"))

	data.PipelineID = pipelineID
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.release.Update(ctx, releaseID, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ReleaseController) Delete(ctx *gin.Context) {
	var releaseID, _ = strconv.Atoi(ctx.Param("releaseID"))

	if err := c.release.Delete(ctx, releaseID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"-1": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ReleaseController) GetLogs(ctx *gin.Context) {
	var releaseID, _ = strconv.Atoi(ctx.Param("releaseID"))

	data, err := c.release.Get(ctx, releaseID)
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
