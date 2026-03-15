package controller

import (
	"net/http"
	"valyria-backend/internal/apps/dragon/service"
	"valyria-backend/internal/pkg/ginhelper"

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
	id, ok := ginhelper.RequireUintParam(ctx, "gitlabID")
	if !ok {
		return
	}
	data, err := c.gitlab.ListGroups(ctx, int64(id))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) ListGroupProjects(ctx *gin.Context) {
	gitlabID, ok := ginhelper.RequireUintParam(ctx, "gitlabID")
	if !ok {
		return
	}
	groupID, ok := ginhelper.RequireUintParam(ctx, "groupID")
	if !ok {
		return
	}
	data, err := c.gitlab.ListGroupProjects(ctx, int64(gitlabID), int64(groupID))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) ListProjectBranches(ctx *gin.Context) {
	gitlabID, ok := ginhelper.RequireUintParam(ctx, "gitlabID")
	if !ok {
		return
	}
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	data, err := c.gitlab.ListProjectBranches(ctx, int64(gitlabID), int64(projectID))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) ListProjectTags(ctx *gin.Context) {
	gitlabID, ok := ginhelper.RequireUintParam(ctx, "gitlabID")
	if !ok {
		return
	}
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	data, err := c.gitlab.ListProjectTags(ctx, int64(gitlabID), int64(projectID))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) GetRefBranch(ctx *gin.Context) {
	gitlabID, ok := ginhelper.RequireUintParam(ctx, "gitlabID")
	if !ok {
		return
	}
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	branch := ctx.Param("branch")
	data, err := c.gitlab.GetRef(ctx, int64(gitlabID), int64(projectID), branch)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *GitlabController) GetRefTag(ctx *gin.Context) {
	gitlabID, ok := ginhelper.RequireUintParam(ctx, "gitlabID")
	if !ok {
		return
	}
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	tag := ctx.Param("tag")
	data, err := c.gitlab.GetRef(ctx, int64(gitlabID), int64(projectID), tag)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
