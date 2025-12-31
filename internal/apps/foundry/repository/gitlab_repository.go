package repository

import (
	"fmt"
	"valyria-backend/internal/apps/foundry/model"

	"gitlab.com/gitlab-org/api/client-go"
)

// GitlabRepository 定义接口
type GitlabRepository interface {
	ListGroups(token, baseURL string) ([]model.GitlabGroup, error)
	ListGroupProjects(token, baseURL string, groupID int) ([]model.GitlabProject, error)
	ListProjectBranches(token, baseURL string, projectID int) ([]model.GitlabRef, error)
	ListProjectTags(token, baseURL string, projectID int) (data []model.GitlabRef, err error)
	GetRef(token, baseURL, ref string, projectID int) (*model.GitlabRef, error)
}

// GitlabRepository 实现了 GitlabRepository 接口
type gitlabRepository struct{}

// NewGitlabRepository 创建新的 GitlabRepository 实例
func NewGitlabRepository() GitlabRepository {
	return &gitlabRepository{}
}

// ListGroups 列表
func (r *gitlabRepository) ListGroups(token, baseURL string) (data []model.GitlabGroup, err error) {
	// 创建 GitLab 客户端
	client, err := GitlabClient(token, baseURL)
	if err != nil {
		return nil, err
	}
	// 调用 GitLab API
	groups, _, err := client.Groups.ListGroups(&gitlab.ListGroupsOptions{})
	for _, g := range groups {
		group := model.GitlabGroup{
			ID:   g.ID,
			Name: g.Name,
			Path: g.Path,
		}
		data = append(data, group)
	}
	return
}

func (r *gitlabRepository) ListGroupProjects(token, baseURL string, groupID int) (data []model.GitlabProject, err error) {
	// 创建 GitLab 客户端
	client, err := GitlabClient(token, baseURL)
	if err != nil {
		return nil, err
	}
	// 调用 GitLab API
	projects, _, err := client.Groups.ListGroupProjects(groupID, &gitlab.ListGroupProjectsOptions{})
	for _, proj := range projects {
		project := model.GitlabProject{
			ID:   proj.ID,
			Name: proj.Name,
			Path: proj.Path,
		}
		data = append(data, project)
	}
	return
}

func (r *gitlabRepository) ListProjectBranches(token, baseURL string, projectID int) (data []model.GitlabRef, err error) {
	// 创建 GitLab 客户端
	client, err := GitlabClient(token, baseURL)
	if err != nil {
		return nil, err
	}
	// 调用 GitLab API
	branches, _, err := client.Branches.ListBranches(projectID, &gitlab.ListBranchesOptions{})
	for _, b := range branches {
		branch := model.GitlabRef{
			RefType:  "branch",
			RefName:  b.Name,
			CommitID: b.Commit.ShortID,
		}
		data = append(data, branch)
	}
	return
}

func (r *gitlabRepository) ListProjectTags(token, baseURL string, projectID int) (data []model.GitlabRef, err error) {
	// 创建 GitLab 客户端
	client, err := GitlabClient(token, baseURL)
	if err != nil {
		return nil, err
	}
	// 调用 GitLab API
	tags, _, err := client.Tags.ListTags(projectID, &gitlab.ListTagsOptions{})
	for _, t := range tags {
		tag := model.GitlabRef{
			RefType:  "tag",
			RefName:  t.Name,
			CommitID: t.Commit.ShortID,
		}
		data = append(data, tag)
	}
	return
}

func (r *gitlabRepository) GetRef(token, baseURL, ref string, projectID int) (data *model.GitlabRef, err error) {
	// 创建 GitLab 客户端
	client, err := GitlabClient(token, baseURL)
	if err != nil {
		return nil, err
	}
	if tag, resp, err := client.Tags.GetTag(projectID, ref); err == nil && resp.StatusCode == 200 {
		return &model.GitlabRef{
			RefType:  "tag",
			RefName:  ref,
			CommitID: tag.Commit.ShortID,
		}, nil
	}

	// 3. 尝试获取 branch
	if branch, resp, err := client.Branches.GetBranch(projectID, ref); err == nil && resp.StatusCode == 200 {
		return &model.GitlabRef{
			RefType:  "branch",
			RefName:  ref,
			CommitID: branch.Commit.ShortID,
		}, nil
	}
	return nil, fmt.Errorf("ref %s not found in project %s", ref, projectID)
}
