package executor

import (
	"fmt"
	"strings"
	"valyria-backend/internal/core/initialize"
	"valyria-backend/internal/core/logger"

	"gitlab.com/gitlab-org/api/client-go"
)

func (p *PipelineExecutor) initGitlabClient() error {
	if p.gitlabClient != nil {
		return nil
	}
	if initialize.Encryptor == nil {
		logger.Errorf("Encryptor is nil")
	}
	token, err := initialize.Encryptor.Decrypt(
		p.release.Pipeline.CredGitlab.Token,
	)
	if err != nil {
		return fmt.Errorf("decrypt gitlab token failed: %w", err)
	}

	client, err := gitlab.NewClient(
		token,
		gitlab.WithBaseURL(p.release.Pipeline.CredGitlab.BaseURL),
	)
	if err != nil {
		return fmt.Errorf("create gitlab client failed: %w", err)
	}
	p.gitlabClient = client
	return nil
}

func (p *PipelineExecutor) resolveGitRepoURL() (repoURL string, err error) {
	fmt.Println("p.release", p.release)
	fmt.Println("p.release.Pipeline ", p.release.Pipeline)
	fmt.Println("p.release.Pipeline.CredGitlab ", p.release.Pipeline.CredGitlab.BaseURL, p.release.Pipeline.CredGitlab.Token)
	if err = p.initGitlabClient(); err != nil {
		return "", err
	}
	project, _, err := p.gitlabClient.Projects.GetProject(p.release.Pipeline.GitlabProjectID, &gitlab.GetProjectOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return repoURL, fmt.Errorf("get gitlab project failed: %w", err)
	}
	fmt.Println(project, project.PathWithNamespace)
	baseURL := strings.TrimRight(p.release.Pipeline.CredGitlab.BaseURL, "/")
	repoURL = fmt.Sprintf("%s/%s.git", baseURL, project.PathWithNamespace)
	return repoURL, nil
}
