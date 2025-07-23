package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	db.AutoMigrate(new(T))
	return &Repository[T]{db: db}
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Model(new(T)).Create(&entity).Error
}

func (r *Repository[T]) Update(ctx context.Context, where T, entity *T) error {
	return r.db.WithContext(ctx).Model(new(T)).Where(where).Updates(&entity).Error
}

func (r *Repository[T]) Delete(ctx context.Context, where T) error {
	return r.db.WithContext(ctx).Model(new(T)).Where(&where).Delete(&where).Error
}

func (r *Repository[T]) Find(ctx context.Context, where T) ([]*T, error) {
	var entities []*T
	err := r.db.WithContext(ctx).Model(new(T)).Where(&where).Find(&entities).Error
	return entities, err
}

func (r *Repository[T]) FindOne(ctx context.Context, where T) (*T, error) {
	var result *T
	err := r.db.WithContext(ctx).Model(new(T)).Where(&where).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository[T]) Count(ctx context.Context, where T) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(new(T)).Where(&where).Count(&count).Error
	return count, err
}

func (r *Repository[T]) Paginate(ctx context.Context, where T, page, pageSize int) ([]T, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	var entities []T
	err := r.db.WithContext(ctx).Model(new(T)).Where(&where).Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities).Error
	return entities, err
}

func (r *Repository[T]) Q(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Model(new(T))
}
