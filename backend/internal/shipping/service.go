package shipping

import (
	"context"
	"sort"
)

type Service struct {
	Repo      Repository
	Providers []ShippingProvider
}

func (s Service) Zones(ctx context.Context) ([]Zone, error) {
	return s.Repo.List(ctx)
}

func (s Service) Quote(ctx context.Context, req QuoteRequest) ([]QuoteOption, error) {
	if req.WeightGrams <= 0 {
		req.WeightGrams = 200
	}
	if req.Dimensions.LengthCM <= 0 {
		req.Dimensions = Dimensions{LengthCM: 20, WidthCM: 12, HeightCM: 8}
	}
	options := []QuoteOption{}
	for _, provider := range s.Providers {
		quoted, err := provider.Quote(ctx, req)
		if err == nil {
			options = append(options, quoted...)
		}
	}
	ownFleet, err := s.Repo.OwnFleetQuote(ctx, req)
	if err != nil {
		return options, err
	}
	options = append(options, ownFleet...)
	sort.SliceStable(options, func(i, j int) bool {
		if options[i].PriceCents == options[j].PriceCents {
			return options[i].ServiceName < options[j].ServiceName
		}
		return options[i].PriceCents < options[j].PriceCents
	})
	return options, nil
}
