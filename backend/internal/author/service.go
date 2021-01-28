package author

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, author Author) error
	GetAllPaginated(ctx context.Context) (*[]Author, error)
	GetByID(ctx context.Context, ID uuid.UUID) (*Author, error)
	UpdateByID(ctx context.Context, ID uuid.UUID) (*Author, error)
	DeleteByID(ctx context.Context, ID uuid.UUID) error
}

type svc struct {
	repo Repository
}

func NewService() Service {
	return &svc{}
}

func (s *svc) Create(ctx context.Context, author Author) error {
	err := s.repo.Insert(ctx, author, nil)
	return err
}

func (s *svc) GetAllPaginated(ctx context.Context) (*[]Author, error) {
	panic("implement me")
}

func (s *svc) GetByID(ctx context.Context, ID uuid.UUID) (*Author, error) {
	panic("implement me")
}

func (s *svc) UpdateByID(ctx context.Context, ID uuid.UUID) (*Author, error) {
	panic("implement me")
}

func (s *svc) DeleteByID(ctx context.Context, ID uuid.UUID) error {
	panic("implement me")
}
