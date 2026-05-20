package shipping

import "context"

type Service struct {
	Repo Repository
}

func (s Service) Zones(ctx context.Context) ([]Zone, error) {
	return s.Repo.List(ctx)
}
