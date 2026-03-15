package executor

import (
	"fmt"

	"gitlab.com/gitlab-org/api/client-go"
)

// GitlabClient 封装 GitLab client 创建
func GitlabClient(token, baseURL string) (*gitlab.Client, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create gitlab client: %w", err)
	}
	// 可在此统一设置超时、分页默认值等
	return client, nil
}
