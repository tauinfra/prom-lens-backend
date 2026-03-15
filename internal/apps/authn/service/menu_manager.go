package service

import (
	"context"
	"sort"
	"valyria-backend/internal/apps/authn/dto"
	"valyria-backend/internal/apps/authn/model"
	"valyria-backend/internal/apps/authn/repository"
	"valyria-backend/internal/apps/authn/request"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type MenuManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.MenuDTO, pg.Pagination, error)
	ListTree(ctx context.Context) ([]dto.MenuTreeDTO, error)
	Get(ctx context.Context, id uint) (dto.MenuDTO, error)
	Create(ctx context.Context, req *request.CreateMenuRequest) error
	Update(ctx context.Context, id uint, req *request.UpdateMenuRequest) error
	Delete(ctx context.Context, id uint) error
}

// menuManager 实现 MenuManager 接口
type menuManager struct {
	repo repository.MenuRepository
	db   *gorm.DB
}

// NewMenuManager 创建新的 MenuManager 实例
func NewMenuManager(db *gorm.DB, repo repository.MenuRepository) MenuManager {
	return &menuManager{db: db, repo: repo}
}

// List 列表
func (s *menuManager) List(ctx context.Context, params pg.QueryParams) (data []dto.MenuDTO, pagination pg.Pagination, err error) {
	menus, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return data, pagination, err
	}
	for _, menu := range menus {
		data = append(data, dto.MenuDTO{
			ID:         menu.ID,
			ParentID:   menu.ParentID,
			Name:       menu.Name,
			Path:       menu.Path,
			Component:  menu.Component,
			Title:      menu.Title,
			Icon:       menu.Icon,
			Rank:       menu.Rank,
			ShowParent: menu.ShowParent,
			ShowLink:   menu.ShowLink,
			Creator:    menu.Creator,
			CreatedAt:  menu.CreatedAt,
			UpdatedAt:  menu.UpdatedAt,
		})
	}
	return data, pagination, nil
}

// ListTree 返回树形菜单（parentId 为空为一级，有 parentId 的为 children）
func (s *menuManager) ListTree(ctx context.Context) ([]dto.MenuTreeDTO, error) {
	menus, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, nil), nil
}

func buildMenuTree(menus []model.Menu, parentID *uint) []dto.MenuTreeDTO {
	var nodes []dto.MenuTreeDTO
	for _, m := range menus {
		if (parentID == nil && m.ParentID == nil) || (parentID != nil && m.ParentID != nil && *m.ParentID == *parentID) {
			node := dto.MenuTreeDTO{
				ID:         m.ID,
				ParentID:   m.ParentID,
				Name:       m.Name,
				Path:       m.Path,
				Component:  m.Component,
				Title:      m.Title,
				Icon:       m.Icon,
				Rank:       m.Rank,
				ShowParent: m.ShowParent,
				ShowLink:   m.ShowLink,
				Creator:    m.Creator,
				CreatedAt:  m.CreatedAt,
				UpdatedAt:  m.UpdatedAt,
			}
			pid := m.ID
			node.Children = buildMenuTree(menus, &pid)
			nodes = append(nodes, node)
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Rank != nodes[j].Rank {
			return nodes[i].Rank < nodes[j].Rank
		}
		return nodes[i].ID < nodes[j].ID
	})
	return nodes
}

// Get 查询
func (s *menuManager) Get(ctx context.Context, id uint) (dto.MenuDTO, error) {
	menu, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.MenuDTO{}, err
	}
	return dto.MenuDTO{
		ID:         menu.ID,
		ParentID:   menu.ParentID,
		Name:       menu.Name,
		Path:       menu.Path,
		Component:  menu.Component,
		Title:      menu.Title,
		Icon:       menu.Icon,
		Rank:       menu.Rank,
		ShowParent: menu.ShowParent,
		ShowLink:   menu.ShowLink,
		Creator:    menu.Creator,
		CreatedAt:  menu.CreatedAt,
		UpdatedAt:  menu.UpdatedAt,
	}, nil
}

// Create 创建
func (s *menuManager) Create(ctx context.Context, req *request.CreateMenuRequest) error {
	data := &model.Menu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Path:       req.Path,
		Component:  req.Component,
		Title:      req.Title,
		Icon:       req.Icon,
		Rank:       req.Rank,
		ShowParent: req.ShowParent,
		ShowLink:   req.ShowLink,
		Creator:    req.Creator,
	}
	return s.repo.Create(ctx, data)
}

// Update 更新
func (s *menuManager) Update(ctx context.Context, id uint, req *request.UpdateMenuRequest) error {
	menu := &model.Menu{}
	menu.ApplyRequest(req)
	return s.repo.Update(ctx, id, menu)
}

// Delete 删除
func (s *menuManager) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
