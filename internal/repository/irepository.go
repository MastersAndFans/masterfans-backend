package repository

import "context"

type IRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	ReadByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, entity *T) error
}
