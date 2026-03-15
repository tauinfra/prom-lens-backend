package service

import (
	"context"
	"errors"
	"sort"
	"time"
	"valyria-backend/internal/apps/authn/dto"
	"valyria-backend/internal/apps/authn/model"
	"valyria-backend/internal/apps/authn/repository"
	"valyria-backend/internal/apps/authn/request"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthManager interface {
	Login(ctx context.Context, username, password string) (model.User, error)
	GetMe(ctx context.Context, id uint) (model.User, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
	ChangePassword(ctx context.Context, id uint, req *request.ChangePasswordRequest) error
	GetRolePermissions(ctx context.Context, userID uint) ([]string, []string, error)
	HasUserPermission(ctx context.Context, userID uint, code string) (bool, error)
	GetUserRoutes(ctx context.Context, userID uint, isSuperuser bool) ([]dto.RouteDTO, error)
}

// authManager 实现 AuthManager 接口
type authManager struct {
	repo repository.AuthRepository
	db   *gorm.DB
}

// NewAuthManager 创建新的 authManager 实例
func NewAuthManager(db *gorm.DB, repo repository.AuthRepository) AuthManager {
	return &authManager{db: db, repo: repo}
}

// Login 登录
func (s *authManager) Login(ctx context.Context, username, password string) (model.User, error) {
	return s.repo.Login(ctx, username, password)
}

// GetMe 获取当前用户信息（用于 profile 等）
func (s *authManager) GetMe(ctx context.Context, id uint) (model.User, error) {
	return s.repo.GetMe(ctx, id)
}

// UpdateLastLogin 登录成功后更新最后登录时间
func (s *authManager) UpdateLastLogin(ctx context.Context, userID uint) error {
	return s.repo.UpdateLastLogin(ctx, userID, time.Now())
}

// ChangePassword 修改密码
func (s *authManager) ChangePassword(ctx context.Context, id uint, req *request.ChangePasswordRequest) error {
	user, err := s.repo.GetMe(ctx, id)
	if err != nil {
		return err
	}
	// 验证旧密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("旧密码错误")
	}
	// 密码新加密
	newPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	// 更新密码
	return s.repo.ChangeMyPassword(ctx, id, string(newPassword))
}

func (s *authManager) GetRolePermissions(ctx context.Context, userID uint) ([]string, []string, error) {
	return s.repo.GetRolePermissions(ctx, userID)
}

func (s *authManager) HasUserPermission(ctx context.Context, userID uint, code string) (bool, error) {
	return s.repo.HasUserPermission(ctx, userID, code)
}

func (s *authManager) GetUserRoutes(ctx context.Context, userID uint, isSuperuser bool) ([]dto.RouteDTO, error) {
	if userID == 0 {
		return nil, errors.New("无效的用户信息")
	}
	var (
		menus []model.Menu
		err   error
	)
	if isSuperuser {
		menus, err = s.repo.ListAllMenus(ctx)
	} else {
		menus, err = s.repo.ListUserMenus(ctx, userID)
	}
	if err != nil {
		return nil, err
	}
	menus, err = s.expandMenuParents(ctx, menus)
	if err != nil {
		return nil, err
	}
	return buildRouteTree(menus), nil
}

func (s *authManager) expandMenuParents(ctx context.Context, menus []model.Menu) ([]model.Menu, error) {
	if len(menus) == 0 {
		return menus, nil
	}
	menuMap := make(map[uint]model.Menu)
	pending := make(map[uint]struct{})
	result := make([]model.Menu, 0, len(menus))
	for _, menu := range menus {
		if _, ok := menuMap[menu.ID]; !ok {
			menuMap[menu.ID] = menu
			result = append(result, menu)
		}
		if menu.ParentID != nil {
			if _, ok := menuMap[*menu.ParentID]; !ok {
				pending[*menu.ParentID] = struct{}{}
			}
		}
	}
	for len(pending) > 0 {
		ids := make([]uint, 0, len(pending))
		for id := range pending {
			ids = append(ids, id)
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		pending = make(map[uint]struct{})
		parents, err := s.repo.GetMenusByIDs(ctx, ids)
		if err != nil {
			return nil, err
		}
		for _, parent := range parents {
			if _, exists := menuMap[parent.ID]; exists {
				continue
			}
			menuMap[parent.ID] = parent
			result = append(result, parent)
			if parent.ParentID != nil {
				if _, ok := menuMap[*parent.ParentID]; !ok {
					pending[*parent.ParentID] = struct{}{}
				}
			}
		}
	}
	return result, nil
}

type routeNode struct {
	value    dto.RouteDTO
	children []*routeNode
}

func buildRouteTree(menus []model.Menu) []dto.RouteDTO {
	nodes := make(map[uint]*routeNode)
	for _, menu := range menus {
		nodes[menu.ID] = &routeNode{
			value: dto.RouteDTO{
				Name:      menu.Name,
				Path:      menu.Path,
				Component: menu.Component,
				Meta: dto.RouteMetaDTO{
					Title:      menu.Title,
					Icon:       menu.Icon,
					Rank:       menu.Rank,
					ShowParent: menu.ShowParent,
					ShowLink:   menu.ShowLink,
				},
			},
		}
	}
	for _, menu := range menus {
		node := nodes[menu.ID]
		if menu.ParentID == nil {
			continue
		}
		parent, ok := nodes[*menu.ParentID]
		if !ok {
			continue
		}
		parent.children = append(parent.children, node)
	}
	roots := make([]*routeNode, 0)
	for _, menu := range menus {
		if menu.ParentID != nil {
			if _, ok := nodes[*menu.ParentID]; ok {
				continue
			}
		}
		if node, ok := nodes[menu.ID]; ok {
			roots = append(roots, node)
		}
	}
	sort.Slice(roots, func(i, j int) bool {
		if roots[i].value.Meta.Rank == roots[j].value.Meta.Rank {
			return roots[i].value.Path < roots[j].value.Path
		}
		return roots[i].value.Meta.Rank < roots[j].value.Meta.Rank
	})
	result := make([]dto.RouteDTO, 0, len(roots))
	for _, node := range roots {
		result = append(result, buildRouteDTO(node))
	}
	return result
}

func buildRouteDTO(node *routeNode) dto.RouteDTO {
	out := node.value
	if len(node.children) > 0 {
		sort.Slice(node.children, func(i, j int) bool {
			if node.children[i].value.Meta.Rank != node.children[j].value.Meta.Rank {
				return node.children[i].value.Meta.Rank < node.children[j].value.Meta.Rank
			}
			return node.children[i].value.Path < node.children[j].value.Path
		})
		out.Children = make([]dto.RouteDTO, 0, len(node.children))
		for _, child := range node.children {
			out.Children = append(out.Children, buildRouteDTO(child))
		}
	}
	return out
}
