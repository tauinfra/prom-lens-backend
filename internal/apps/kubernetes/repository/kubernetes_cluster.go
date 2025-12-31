package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"valyria-backend/internal/apps/kubernetes/model"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/encryption"
)

// ClusterRepository 定义接口
type ClusterRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Cluster, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Cluster, error)
	Create(ctx context.Context, data *model.Cluster) error
	Update(ctx context.Context, id int, cluster *model.Cluster) error
	Delete(ctx context.Context, id int) error
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
func (r *clusterRepository) Get(ctx context.Context, id int) (model.Cluster, error) {
	var data model.Cluster
	if err := r.db.WithContext(ctx).First(&data, id).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *clusterRepository) Create(ctx context.Context, data *model.Cluster) error {
	// 创建前加密 Token
	encryptToken, err := r.encryptor.Encrypt(data.Token)
	if err != nil {
		return fmt.Errorf("kubernetes encrypt token failed, err: %v", err)
	}
	data.Token = encryptToken
	if err := r.db.WithContext(ctx).Model(&model.Cluster{}).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *clusterRepository) Update(ctx context.Context, id int, data *model.Cluster) error {
	// 1. 先查询现有数据
	var cluster model.Cluster
	if err := r.db.WithContext(ctx).First(&cluster, id).Error; err != nil {
		return err
	}
	// 2. 检查 Token 是否有变化
	encryptTokenChanged := data.Token != "" && data.Token != cluster.Token
	if encryptTokenChanged {
		encryptToken, err := r.encryptor.Encrypt(data.Token)
		if err != nil {
			return fmt.Errorf("kubernetes encrypt token failed, err: %v", err)
		}
		data.Token = encryptToken
	}
	// 3. 更新集群
	if err := r.db.WithContext(ctx).Model(&model.Cluster{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	// 4. 查询最新数据
	return r.db.WithContext(ctx).First(&data, id).Error
}

// Delete 删除
func (r *clusterRepository) Delete(ctx context.Context, id int) error {
	if err := r.db.WithContext(ctx).Delete(&model.Cluster{}, id).Error; err != nil {
		return err
	}
	return nil
}

// WithTx 返回一个绑定事务的 Repository
func (r *clusterRepository) WithTx(db *gorm.DB) ClusterRepository {
	return &clusterRepository{db: db}
}
