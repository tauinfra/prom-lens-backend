package executor

type ReleasePhase string
type ReleaseStatus string

const (
	PhaseStarting  ReleasePhase = "starting"
	PhaseCompleted ReleasePhase = "completed"
)

const (
	StatusPending ReleaseStatus = "pending"
	StatusRunning ReleaseStatus = "progressing"
	StatusSuccess ReleaseStatus = "success"
	StatusFailed  ReleaseStatus = "failed"
)
