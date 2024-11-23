package service

import "context"

type repository interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type DomainService struct {
	repo repository
}

func New(repo repository) *DomainService {
	s := &DomainService{
		repo: repo,
	}
	return s
}

func (s *DomainService) Run() {
}
