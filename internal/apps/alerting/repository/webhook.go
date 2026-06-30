package repository

import (
	"context"

	"prom-lens-backend/internal/apps/alerting/model"
	pg "prom-lens-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type WebhookRepository interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Webhook, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Webhook, error)
	GetByName(ctx context.Context, name string) (model.Webhook, error)
	ListAllWithRoutes(ctx context.Context) ([]model.Webhook, error)
	Create(ctx context.Context, data *model.Webhook, route *model.Route, matchers []model.RouteMatcher) error
	Update(ctx context.Context, id int, data *model.Webhook, route *model.Route, matchers []model.RouteMatcher, replaceRoute bool) error
	Delete(ctx context.Context, id int) error
}

type webhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) WebhookRepository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) preloadRoute(db *gorm.DB) *gorm.DB {
	return db.Preload("Route.Matchers", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sort_order ASC, id ASC")
	})
}

func (r *webhookRepository) List(ctx context.Context, params pg.QueryParams) ([]model.Webhook, pg.Pagination, error) {
	var data []model.Webhook
	q := r.preloadRoute(r.db.WithContext(ctx))
	pagination, err := pg.Paginate(q, &data, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	return data, pagination, nil
}

func (r *webhookRepository) Get(ctx context.Context, id int) (model.Webhook, error) {
	var data model.Webhook
	if err := r.preloadRoute(r.db.WithContext(ctx)).First(&data, id).Error; err != nil {
		return model.Webhook{}, err
	}
	return data, nil
}

func (r *webhookRepository) GetByName(ctx context.Context, name string) (model.Webhook, error) {
	var data model.Webhook
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&data).Error; err != nil {
		return model.Webhook{}, err
	}
	return data, nil
}

func (r *webhookRepository) ListAllWithRoutes(ctx context.Context) ([]model.Webhook, error) {
	var data []model.Webhook
	if err := r.preloadRoute(r.db.WithContext(ctx)).Order("id ASC").Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *webhookRepository) Create(ctx context.Context, data *model.Webhook, route *model.Route, matchers []model.RouteMatcher) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(data).Error; err != nil {
			return err
		}
		if route == nil || len(matchers) == 0 {
			return nil
		}
		route.WebhookID = data.ID
		if err := tx.Create(route).Error; err != nil {
			return err
		}
		for i := range matchers {
			matchers[i].RouteID = route.ID
			matchers[i].SortOrder = i
			if err := tx.Create(&matchers[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *webhookRepository) Update(ctx context.Context, id int, data *model.Webhook, route *model.Route, matchers []model.RouteMatcher, replaceRoute bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Updates(data).Error; err != nil {
			return err
		}
		if !replaceRoute {
			return nil
		}
		if err := tx.Where("webhook_id = ?", id).Delete(&model.Route{}).Error; err != nil {
			return err
		}
		if route == nil || len(matchers) == 0 {
			return nil
		}
		route.WebhookID = id
		if err := tx.Create(route).Error; err != nil {
			return err
		}
		for i := range matchers {
			matchers[i].RouteID = route.ID
			matchers[i].SortOrder = i
			if err := tx.Create(&matchers[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *webhookRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Webhook{}).Error
}
