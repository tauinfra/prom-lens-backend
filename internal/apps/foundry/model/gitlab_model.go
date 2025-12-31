package model

// GitlabRepoList 代表整个 GitLab 仓库列表
type GitlabRepoList struct {
	Groups []GitlabGroup `json:"groups"` // 顶层 group 列表
}

// GitlabGroup 代表一个 GitLab group
type GitlabGroup struct {
	ID       int64           `json:"id,omitempty"`
	Name     string          `json:"name,omitempty"`
	Path     string          `json:"path,omitempty"`     // group 路径
	Projects []GitlabProject `json:"projects,omitempty"` // group 下的项目列表
}

// GitlabProject 代表一个 GitLab 项目
type GitlabProject struct {
	ID       int64       `json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`     // 项目名
	Path     string      `json:"path,omitempty"`     // 项目路径
	Branches []GitlabRef `json:"branches,omitempty"` // 项目下的分支列表
}

// GitlabBranch 代表一个 GitLab 分支

type GitlabRef struct {
	RefType  string
	RefName  string
	CommitID string
}
