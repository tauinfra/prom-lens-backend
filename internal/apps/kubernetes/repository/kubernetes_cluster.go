package repository

import (
	"context"
	"fmt"
	"valyria-backend/internal/apps/kubernetes/model"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

// ClusterRepository 定义接口
type ClusterRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Cluster, pg.Pagination, error)
	Get(ctx context.Context, id uint) (model.Cluster, error)
	Create(ctx context.Context, data *model.Cluster) error
	Update(ctx context.Context, id uint, cluster *model.Cluster) error
	UpdateToken(ctx context.Context, id uint, token string) error
	Delete(ctx context.Context, id uint) error
	WithTx(tx *gorm.DB) ClusterRepository // 提供一个事务操作接口

}

// clusterRepository 实现了 ClusterRepository 接口
type clusterRepository struct {
	db        *gorm.DB
	encryptor encryption.Encryptor
}

// NewClusterRepository 创建新的 ClusterRepository 实例
func NewClusterRepository(db *gorm.DB, encryptor encryption.Encryptor) ClusterRepository {
	return &clusterRepository{db: db, encryptor: encryptor}
}

// List 列表
func (r *clusterRepository) List(ctx context.Context, params pg.QueryParams) (data []model.Cluster, pagination pg.Pagination, err error) {
	// 分页查询
	if pagination, err = pg.Paginate(r.db.WithContext(ctx), &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *clusterRepository) Get(ctx context.Context, id uint) (data model.Cluster, err error) {
	err = r.db.WithContext(ctx).First(&data, id).Error
	return
}

// Create 创建
func (r *clusterRepository) Create(ctx context.Context, data *model.Cluster) error {
	// 创建前加密 Token
	encryptToken, err := r.encryptor.Encrypt(data.Token)
	if err != nil {
		return fmt.Errorf("kubernetes encrypt token failed, err: %v", err)
	}
	data.Token = encryptToken
	return r.db.WithContext(ctx).Model(&model.Cluster{}).Create(data).Error
}

// Update 更新
func (r *clusterRepository) Update(ctx context.Context, id uint, data *model.Cluster) error {
	// 更新时排除 token 字段，禁止修改 token
	// 使用 Omit 排除 token 字段
	return r.db.WithContext(ctx).Model(&model.Cluster{}).
		Where("id = ?", id).
		Omit("token").
		Updates(data).Error
}

// UpdateToken 更新集群 token（入库前加密）
func (r *clusterRepository) UpdateToken(ctx context.Context, id uint, token string) error {
	encryptToken, err := r.encryptor.Encrypt(token)
	if err != nil {
		return fmt.Errorf("kubernetes encrypt token failed, err: %v", err)
	}
	return r.db.WithContext(ctx).Model(&model.Cluster{}).
		Where("id = ?", id).
		Update("token", encryptToken).Error
}

// Delete 删除
func (r *clusterRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Cluster{}, id).Error; err != nil {
		return err
	}
	return nil
}

// WithTx 返回一个绑定事务的 Repository
func (r *clusterRepository) WithTx(db *gorm.DB) ClusterRepository {
	return &clusterRepository{db: db, encryptor: r.encryptor}
}
