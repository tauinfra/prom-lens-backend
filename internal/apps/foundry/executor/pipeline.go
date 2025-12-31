package executor

import (
	"fmt"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *PipelineExecutor) pipelineRunTemplate() *tektonv1.PipelineRun {
	kustomizePath := fmt.Sprintf(
		"kustomize/%s/%s/overlays/%s",
		p.release.Pipeline.Environment.Project.Name,
		p.release.Pipeline.Application.Name,
		p.release.Pipeline.Environment.Name,
	)
	app := fmt.Sprintf("%s-%s-%s",
		p.release.Pipeline.Environment.Project.Name,
		p.release.Pipeline.Environment.Name,
		p.release.Pipeline.Application.Name,
	)
	repoURL, _ := p.resolveGitRepoURL()
	return &tektonv1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.release.TaskID,
			Namespace: "default",
		},
		Spec: tektonv1.PipelineRunSpec{
			PipelineRef: &tektonv1.PipelineRef{
				Name: p.release.Pipeline.Template,
			},
			Workspaces: []tektonv1.WorkspaceBinding{
				{
					Name:                "workspace",
					VolumeClaimTemplate: nil,
					EmptyDir:            &corev1.EmptyDirVolumeSource{},
				},
			},
			Params: []tektonv1.Param{
				stringParam("image", "nginx"),
				stringParam("image_tag", "1.24"),
				stringParam("code_repo_url", repoURL),
				stringParam("code_repo_revision", p.release.GitRef),
				stringParam("config_repo_url", "https://gitlab.com/valyrias/kubernetes.git"),
				stringParam("config_repo_revision", "main"),
				stringParam("kustomize_path", kustomizePath),
				stringParam("argocd_server", "argocd-server.argocd.svc"),
				stringParam("argocd_app", app),
				stringParam("script", p.release.Pipeline.Script),
			},
		},
	}
}

func stringParam(name, value string) tektonv1.Param {
	return tektonv1.Param{
		Name: name,
		Value: tektonv1.ParamValue{
			Type:      tektonv1.ParamTypeString,
			StringVal: value,
		},
	}
}
