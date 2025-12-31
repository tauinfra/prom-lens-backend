package executor

import (
	"sync"

	"valyria-backend/internal/apps/foundry/model"

	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"gitlab.com/gitlab-org/api/client-go"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
)

type PipelineExecutor struct {
	gitlabClient  *gitlab.Client
	tektonClient  *tektonclient.Clientset
	kubeClient    *kubernetes.Clientset
	db            *gorm.DB
	release       model.Release
	taskOrder     []string
	status        ReleaseStatus
	mu            sync.Mutex
	log           *logWriter
	processedPods sync.Map // key = pod/container
}
