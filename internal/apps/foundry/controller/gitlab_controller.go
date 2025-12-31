package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/foundry/service"

	"github.com/gin-gonic/gin"
)

// GitlabController 定义控制器结构体
type GitlabController struct {
	gitlab service.GitlabManager // 使用服务接口
}

func NewGitlabController(gitlab service.GitlabManager) *GitlabController {
	return &GitlabController{gitlab: gitlab}
}

func (c *GitlabController) ListGitlabGroups(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("gitlabID"))
	data, err := c.gitlab.ListGroups(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) ListGroupProjects(ctx *gin.Context) {
	gitlabID, _ := strconv.Atoi(ctx.Param("gitlabID"))
	groupID, _ := strconv.Atoi(ctx.Param("groupID"))
	data, err := c.gitlab.ListGroupProjects(ctx, gitlabID, groupID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) ListProjectBranches(ctx *gin.Context) {
	gitlabID, _ := strconv.Atoi(ctx.Param("gitlabID"))
	projectID, _ := strconv.Atoi(ctx.Param("projectID"))
	data, err := c.gitlab.ListProjectBranches(ctx, gitlabID, projectID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) ListProjectTags(ctx *gin.Context) {
	gitlabID, _ := strconv.Atoi(ctx.Param("gitlabID"))
	projectID, _ := strconv.Atoi(ctx.Param("projectID"))
	data, err := c.gitlab.ListProjectTags(ctx, gitlabID, projectID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) GetRefBranch(ctx *gin.Context) {
	gitlabID, _ := strconv.Atoi(ctx.Param("gitlabID"))
	projectID, _ := strconv.Atoi(ctx.Param("projectID"))
	branch := ctx.Param("branch")
	data, err := c.gitlab.GetRef(ctx, gitlabID, projectID, branch)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) GetRefTag(ctx *gin.Context) {
	gitlabID, _ := strconv.Atoi(ctx.Param("gitlabID"))
	projectID, _ := strconv.Atoi(ctx.Param("projectID"))
	tag := ctx.Param("tag")
	data, err := c.gitlab.GetRef(ctx, gitlabID, projectID, tag)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
