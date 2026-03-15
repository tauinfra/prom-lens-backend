package service

import (
	"context"
	"fmt"
	"time"
	domainpermission "valyria-backend/internal/apps/kubernetes/domain/permission"
	"valyria-backend/internal/apps/kubernetes/dto"
	"valyria-backend/internal/apps/kubernetes/model"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/request"
	"valyria-backend/internal/core/logger"
)

const (
	maxLastErrorLen = 512
	defaultMaxRetry = 5 // 与 config permission_worker 默认一致，仅当未注入 policy 时使用
)

func truncateErr(s string) string {
	if len(s) <= maxLastErrorLen {
		return s
	}
	return s[:maxLastErrorLen-3] + "..."
}

// SyncPolicy 同步重试策略（由 config permission_worker 解析后注入）
type SyncPolicy struct {
	MaxRetry  int           // 最大重试次数，达此后置为 dead
	BaseDelay time.Duration // 退避基准（base_delay）
	MaxDelay  time.Duration // 退避上限（max_delay）
	BatchSize int           // 每批拉取条数
}

// nextRetryAt 退避：baseDelay * 2^retryCount，上限 maxDelay；若 retryCount >= maxRetry 返回 nil（标记 dead）
func nextRetryAt(retryCount int, maxRetry int, baseDelay, maxDelay time.Duration) *time.Time {
	if maxRetry <= 0 {
		maxRetry = defaultMaxRetry
	}
	if retryCount >= maxRetry {
		return nil
	}
	delay := baseDelay
	if delay <= 0 {
		delay = time.Minute
	}
	delay = delay * time.Duration(1<<uint(retryCount))
	if maxDelay > 0 && delay > maxDelay {
		delay = maxDelay
	}
	t := time.Now().Add(delay)
	return &t
}

type PermissionManager interface {
	List(ctx context.Context, clusterID uint) ([]dto.PermissionDTO, error)
	Create(ctx context.Context, req *request.CreatePermissionRequest, creator string) error
	CreateBatch(ctx context.Context, req *request.BatchCreatePermissionRequest, creator string) error
	Delete(ctx context.Context, id uint) error
	SyncPendingPermissions(ctx context.Context) error
}

type permissionManager struct {
	repo                   repository.PermissionRepository
	clusterRoleBindingRepo repository.ClusterRoleBindingRepository
	roleBindingRepo        repository.RoleBindingRepository
	policy                 SyncPolicy
}

func NewPermissionManager(repo repository.PermissionRepository, clusterRoleBindingRepo repository.ClusterRoleBindingRepository, roleBindingRepo repository.RoleBindingRepository, policy SyncPolicy) PermissionManager {
	if policy.BatchSize <= 0 {
		policy.BatchSize = 100
	}
	if policy.MaxRetry <= 0 {
		policy.MaxRetry = defaultMaxRetry
	}
	return &permissionManager{
		repo:                   repo,
		clusterRoleBindingRepo: clusterRoleBindingRepo,
		roleBindingRepo:        roleBindingRepo,
		policy:                 policy,
	}
}

func (s *permissionManager) List(ctx context.Context, clusterID uint) ([]dto.PermissionDTO, error) {
	list, err := s.repo.List(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.PermissionDTO, 0, len(list))
	for _, m := range list {
		username, clusterName := "", ""
		if m.User != nil {
			username = m.User.Username
		}
		if m.Cluster != nil {
			clusterName = m.Cluster.Name
		}
		syncStatus := m.SyncStatus
		if syncStatus == "" {
			syncStatus = model.SyncStatusSynced
		}
		out = append(out, dto.PermissionDTO{
			ID:          m.ID,
			UserID:      m.UserID,
			Username:    username,
			ClusterID:   m.ClusterID,
			ClusterName: clusterName,
			Namespace:   m.Namespace,
			Role:        m.Role,
			Creator:     m.Creator,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
			SyncStatus:  syncStatus,
			RetryCount:  m.RetryCount,
			LastError:   m.LastError,
			LastErrorAt: m.LastErrorAt,
			LastSyncAt:  m.LastSyncAt,
			NextRetryAt: m.NextRetryAt,
		})
	}
	return out, nil
}

func (s *permissionManager) Create(ctx context.Context, req *request.CreatePermissionRequest, creator string) error {
	if !domainpermission.IsValidRole(req.Role) {
		return domainpermission.ErrInvalidRole
	}
	username, err := s.repo.GetUsernameByUserID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("get username for userID %d: %w", req.UserID, err)
	}
	if username == "" {
		return fmt.Errorf("user not found or username empty (userID=%d), cannot sync K8s RBAC", req.UserID)
	}
	m := &model.Permission{
		UserID:     req.UserID,
		ClusterID:  req.ClusterID,
		Namespace:  req.Namespace,
		Role:       req.Role,
		Creator:    creator,
		SyncStatus: model.SyncStatusPending,
	}
	if err := s.repo.Create(ctx, m); err != nil {
		return err
	}
	// 尝试同步 K8s；失败则标记 failed，由后台任务重试，接口仍返回成功
	if IsClusterScoped(m.Namespace) {
		logger.Infof("permission create: creating ClusterRoleBinding clusterID=%d user=%s role=%s", int(m.ClusterID), username, m.Role)
		if err := SyncClusterRoleBindingCreate(ctx, uint(m.ClusterID), username, m.Role, s.clusterRoleBindingRepo.Create); err != nil {
			logger.Errorf("permission create: ClusterRoleBinding create failed: %v", err)
			_ = s.repo.UpdateSyncFailure(ctx, m.ID, 1, truncateErr(err.Error()), nextRetryAt(1, s.policy.MaxRetry, s.policy.BaseDelay, s.policy.MaxDelay), s.policy.MaxRetry)
			return nil
		}
		now := time.Now()
		_ = s.repo.UpdateSyncSuccess(ctx, m.ID, now)
		logger.Infof("permission create: ClusterRoleBinding created clusterID=%d user=%s role=%s", int(m.ClusterID), username, m.Role)
		return nil
	}
	logger.Infof("permission create: creating RoleBinding clusterID=%d ns=%s user=%s role=%s", int(m.ClusterID), m.Namespace, username, m.Role)
	if err := SyncRoleBindingCreate(ctx, uint(m.ClusterID), m.Namespace, username, m.Role, s.roleBindingRepo.Create); err != nil {
		logger.Errorf("permission create: RoleBinding create failed: %v", err)
		_ = s.repo.UpdateSyncFailure(ctx, m.ID, 1, truncateErr(err.Error()), nextRetryAt(1, s.policy.MaxRetry, s.policy.BaseDelay, s.policy.MaxDelay), s.policy.MaxRetry)
		return nil
	}
	now := time.Now()
	_ = s.repo.UpdateSyncSuccess(ctx, m.ID, now)
	logger.Infof("permission create: RoleBinding created clusterID=%d ns=%s user=%s role=%s", int(m.ClusterID), m.Namespace, username, m.Role)
	return nil
}

// permKey 用于 diff 的 (ClusterID, Namespace, Role) 唯一键
type permKey struct {
	ClusterID uint
	Namespace string
	Role      string
}

func (s *permissionManager) CreateBatch(ctx context.Context, req *request.BatchCreatePermissionRequest, creator string) error {
	// 1. 获取该用户当前权限
	current, err := s.repo.ListByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}
	// 2. 从请求构建目标集合（按 cluster 分组展开，去重）
	desiredSet := make(map[permKey]struct{})
	var desiredList []permKey
	for _, c := range req.Clusters {
		for _, b := range c.Bindings {
			if !domainpermission.IsValidRole(b.Role) {
				return domainpermission.ErrInvalidRole
			}
			k := permKey{ClusterID: c.ClusterID, Namespace: b.Namespace, Role: b.Role}
			if _, ok := desiredSet[k]; !ok {
				desiredSet[k] = struct{}{}
				desiredList = append(desiredList, k)
			}
		}
	}
	currentSet := make(map[permKey]*model.Permission)
	for i := range current {
		p := &current[i]
		k := permKey{ClusterID: p.ClusterID, Namespace: p.Namespace, Role: p.Role}
		currentSet[k] = p
	}
	var toRemove []*model.Permission
	for k, p := range currentSet {
		if _, ok := desiredSet[k]; !ok {
			toRemove = append(toRemove, p)
		}
	}
	var toAdd []permKey
	for _, k := range desiredList {
		if _, ok := currentSet[k]; !ok {
			toAdd = append(toAdd, k)
		}
	}
	logger.Infof("permission batch: userID=%d toRemove=%d toAdd=%d", req.UserID, len(toRemove), len(toAdd))
	// 3. 删除 RBAC
	for _, p := range toRemove {
		username := ""
		if p.User != nil {
			username = p.User.Username
		}
		if IsClusterScoped(p.Namespace) {
			logger.Infof("permission batch: removing ClusterRoleBinding clusterID=%d user=%s role=%s", int(p.ClusterID), username, p.Role)
			_ = SyncClusterRoleBindingDelete(ctx, uint(p.ClusterID), username, p.Role, s.clusterRoleBindingRepo.Delete)
		} else {
			logger.Infof("permission batch: removing RoleBinding clusterID=%d ns=%s user=%s role=%s", int(p.ClusterID), p.Namespace, username, p.Role)
			_ = SyncRoleBindingDelete(ctx, uint(p.ClusterID), p.Namespace, username, p.Role, s.roleBindingRepo.Delete)
		}
	}
	// 4. 删除 DB
	for _, p := range toRemove {
		logger.Infof("permission batch: removing permission id=%d clusterID=%d ns=%s role=%s", p.ID, int(p.ClusterID), p.Namespace, p.Role)
		if err := s.repo.Delete(ctx, p.ID); err != nil {
			return err
		}
	}
	// 若有新增项，先解析用户名（用于 K8s Subject）
	var username string
	if len(toAdd) > 0 {
		username, err = s.repo.GetUsernameByUserID(ctx, req.UserID)
		if err != nil {
			return fmt.Errorf("get username for userID %d: %w", req.UserID, err)
		}
		if username == "" {
			return fmt.Errorf("user not found or username empty (userID=%d), cannot sync K8s RBAC", req.UserID)
		}
	}
	// 5. 创建 DB（pending_sync），6. 尝试 K8s 同步；失败由后台任务重试，不失败整批请求
	for _, k := range toAdd {
		m := &model.Permission{
			UserID:     req.UserID,
			ClusterID:  k.ClusterID,
			Namespace:  k.Namespace,
			Role:       k.Role,
			Creator:    creator,
			SyncStatus: model.SyncStatusPending,
		}
		if err := s.repo.Create(ctx, m); err != nil {
			return err
		}
		if IsClusterScoped(k.Namespace) {
			logger.Infof("permission batch: creating ClusterRoleBinding clusterID=%d user=%s role=%s", int(k.ClusterID), username, k.Role)
			if err := SyncClusterRoleBindingCreate(ctx, uint(k.ClusterID), username, k.Role, s.clusterRoleBindingRepo.Create); err != nil {
				logger.Errorf("permission batch: ClusterRoleBinding create failed: %v", err)
				_ = s.repo.UpdateSyncFailure(ctx, m.ID, 1, truncateErr(err.Error()), nextRetryAt(1, s.policy.MaxRetry, s.policy.BaseDelay, s.policy.MaxDelay), s.policy.MaxRetry)
			} else {
				_ = s.repo.UpdateSyncSuccess(ctx, m.ID, time.Now())
			}
		} else {
			logger.Infof("permission batch: creating RoleBinding clusterID=%d ns=%s user=%s role=%s", int(k.ClusterID), k.Namespace, username, k.Role)
			if err := SyncRoleBindingCreate(ctx, uint(k.ClusterID), k.Namespace, username, k.Role, s.roleBindingRepo.Create); err != nil {
				logger.Errorf("permission batch: RoleBinding create failed: %v", err)
				_ = s.repo.UpdateSyncFailure(ctx, m.ID, 1, truncateErr(err.Error()), nextRetryAt(1, s.policy.MaxRetry, s.policy.BaseDelay, s.policy.MaxDelay), s.policy.MaxRetry)
			} else {
				_ = s.repo.UpdateSyncSuccess(ctx, m.ID, time.Now())
			}
		}
	}
	return nil
}

func (s *permissionManager) Delete(ctx context.Context, id uint) error {
	perm, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	username := ""
	if perm.User != nil {
		username = perm.User.Username
	} else {
		username, _ = s.repo.GetUsernameByUserID(ctx, perm.UserID)
	}
	if IsClusterScoped(perm.Namespace) {
		_ = SyncClusterRoleBindingDelete(ctx, uint(perm.ClusterID), username, perm.Role, s.clusterRoleBindingRepo.Delete)
	} else {
		_ = SyncRoleBindingDelete(ctx, uint(perm.ClusterID), perm.Namespace, username, perm.Role, s.roleBindingRepo.Delete)
	}
	return nil
}

// SyncPendingPermissions 后台任务：拉取待同步/失败记录并重试 K8s 同步（使用 policy.BatchSize、policy.MaxRetry）
func (s *permissionManager) SyncPendingPermissions(ctx context.Context) error {
	limit := s.policy.BatchSize
	if limit <= 0 {
		limit = 100
	}
	list, err := s.repo.ListPendingSync(ctx, limit, s.policy.MaxRetry)
	if err != nil {
		return err
	}
	for i := range list {
		p := &list[i]
		username, err := s.repo.GetUsernameByUserID(ctx, p.UserID)
		if err != nil || username == "" {
			next := p.RetryCount + 1
			_ = s.repo.UpdateSyncFailure(ctx, p.ID, next, "get username failed or empty", nextRetryAt(next, s.policy.MaxRetry, s.policy.BaseDelay, s.policy.MaxDelay), s.policy.MaxRetry)
			continue
		}
		if IsClusterScoped(p.Namespace) {
			if err := SyncClusterRoleBindingCreate(ctx, uint(p.ClusterID), username, p.Role, s.clusterRoleBindingRepo.Create); err != nil {
				next := p.RetryCount + 1
				logger.Warnf("permission sync: ClusterRoleBinding retry failed id=%d: %v", p.ID, err)
				_ = s.repo.UpdateSyncFailure(ctx, p.ID, next, truncateErr(err.Error()), nextRetryAt(next, s.policy.MaxRetry, s.policy.BaseDelay, s.policy.MaxDelay), s.policy.MaxRetry)
			} else {
				_ = s.repo.UpdateSyncSuccess(ctx, p.ID, time.Now())
				logger.Infof("permission sync: ClusterRoleBinding synced id=%d", p.ID)
			}
		} else {
			if err := SyncRoleBindingCreate(ctx, uint(p.ClusterID), p.Namespace, username, p.Role, s.roleBindingRepo.Create); err != nil {
				next := p.RetryCount + 1
				logger.Warnf("permission sync: RoleBinding retry failed id=%d: %v", p.ID, err)
				_ = s.repo.UpdateSyncFailure(ctx, p.ID, next, truncateErr(err.Error()), nextRetryAt(next, s.policy.MaxRetry, s.policy.BaseDelay, s.policy.MaxDelay), s.policy.MaxRetry)
			} else {
				_ = s.repo.UpdateSyncSuccess(ctx, p.ID, time.Now())
				logger.Infof("permission sync: RoleBinding synced id=%d", p.ID)
			}
		}
	}
	return nil
}
