package service

import "context"

type AMSyncResult struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg,omitempty"`
}

func runAMSync(syncer interface {
	SyncAll(ctx context.Context) error
}, ctx context.Context) AMSyncResult {
	if syncer == nil {
		return AMSyncResult{Success: true, Msg: "alertmanager sync skipped (not configured)"}
	}
	if err := syncer.SyncAll(ctx); err != nil {
		return AMSyncResult{Success: false, Msg: err.Error()}
	}
	return AMSyncResult{Success: true}
}
