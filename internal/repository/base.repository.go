package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("record not found")

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) FindById(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &entity, err
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, limit, offset int) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var entity T
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
}

func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).Count(&count).Error
	return count, err
}
