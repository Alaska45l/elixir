package shipping

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"sort"
	"strings"
)

var postalCodeRE = regexp.MustCompile(`^[A-Z]?\d{4}[A-Z]{0,3}$`)

type Service struct {
	Repo      Repository
	Providers []ShippingProvider
}

func (s Service) Zones(ctx context.Context) ([]Zone, error) {
	return s.Repo.List(ctx)
}

func (s Service) Quote(ctx context.Context, req QuoteRequest) ([]QuoteOption, error) {
	if err := normalizeQuoteRequest(&req); err != nil {
		return nil, err
	}
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
		} else if !errors.Is(err, ErrProviderUnavailable) {
			slog.Warn("shipping provider unavailable", "error", err)
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

func normalizeQuoteRequest(req *QuoteRequest) error {
	req.DestinationPostalCode = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(req.DestinationPostalCode), " ", ""))
	req.ProvinceCode = strings.ToUpper(strings.TrimSpace(req.ProvinceCode))
	if !postalCodeRE.MatchString(req.DestinationPostalCode) {
		return fmt.Errorf("%w: código postal inválido", ErrInvalidQuoteRequest)
	}
	if len(req.ProvinceCode) > 40 {
		return fmt.Errorf("%w: provincia inválida", ErrInvalidQuoteRequest)
	}
	if req.WeightGrams < 0 || req.WeightGrams > 50000 {
		return fmt.Errorf("%w: peso inválido", ErrInvalidQuoteRequest)
	}
	if req.Dimensions.LengthCM < 0 || req.Dimensions.WidthCM < 0 || req.Dimensions.HeightCM < 0 {
		return fmt.Errorf("%w: dimensiones inválidas", ErrInvalidQuoteRequest)
	}
	if req.Dimensions.LengthCM > 200 || req.Dimensions.WidthCM > 200 || req.Dimensions.HeightCM > 200 {
		return fmt.Errorf("%w: dimensiones inválidas", ErrInvalidQuoteRequest)
	}
	return nil
}
