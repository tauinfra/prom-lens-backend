package service

const (
	ActionView     = "view"
	ActionDeploy   = "deploy"
	ActionApprove  = "approve"
	ActionRollback = "rollback"
)

func IsValidAction(action string) bool {
	switch action {
	case ActionView, ActionDeploy, ActionApprove, ActionRollback:
		return true
	default:
		return false
	}
}
