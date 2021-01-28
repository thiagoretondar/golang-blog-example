package database

import (
	"context"
)

type CRUDRepository interface {
	Insert(ctx context.Context, data interface{}, outputInsertedID interface{}) error
	Find(ctx context.Context, filter interface{}, output interface{}) error
	FindAll(ctx context.Context, output interface{}) error
	FindOne(ctx context.Context, filter interface{}, result interface{}) error
	Update(ctx context.Context, set map[string]interface{}, filter interface{}) (int64, error)
	Remove(ctx context.Context, filter interface{}, physicalDeletion bool) (int64, error)
	Count(ctx context.Context, filter interface{}) (int64, error)
	ExtractColumnPairs(data interface{}) ([]string, []interface{})
}
