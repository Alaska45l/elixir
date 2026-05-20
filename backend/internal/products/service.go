package products

import "context"

type Service struct {
	Repo Repository
}

func (s Service) List(ctx context.Context, f ListFilters) ([]Product, error) {
	return s.Repo.List(ctx, f)
}

func (s Service) Detail(ctx context.Context, slug string) (*Product, error) {
	return s.Repo.BySlug(ctx, slug)
}

func (s Service) Search(ctx context.Context, q string) ([]SearchResult, error) {
	return s.Repo.Search(ctx, q)
}
