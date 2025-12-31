package repository

import (
	"fmt"

	"gitlab.com/gitlab-org/api/client-go"
)

type ClientFactory struct{}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{}
}

// Create 使用凭证创建 GitLab 客户端
func (f *ClientFactory) Create(token, baseURL string) (git *gitlab.Client, err error) {
	git, err = gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to create gitlab client with BaseURL: %v failed. err: %w", baseURL, err)
	}
	return
}
