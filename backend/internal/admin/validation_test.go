package admin

import (
	"strings"
	"testing"
	"time"
)

func TestNormalizeProductPayloadDefaultsAndValidates(t *testing.T) {
	p := productPayload{
		Name:          " Nocturno Oud ",
		Slug:          "nocturno-oud",
		ScentFamily:   "Amaderado",
		GenderTag:     "Unisex",
		Concentration: "EDP",
		TopNotes:      []string{"Bergamota", "Bergamota", strings.Repeat("x", 50)},
		Variants:      []variantForm{{SizeML: 50, PriceARSCents: 8900000, Stock: 4}},
		Images:        []imageForm{{URL: "https://example.com/image.jpg", AltText: "", IsPrimary: false}},
	}

	if err := normalizeProductPayload(&p); err != nil {
		t.Fatalf("expected valid product: %v", err)
	}
	if p.Name != "Nocturno Oud" {
		t.Fatalf("name was not trimmed: %q", p.Name)
	}
	if p.Variants[0].SKU != "nocturno-oud-50" {
		t.Fatalf("expected generated sku, got %q", p.Variants[0].SKU)
	}
	if !p.Images[0].IsPrimary || p.Images[0].AltText != "Nocturno Oud" {
		t.Fatalf("expected primary image defaults, got %#v", p.Images[0])
	}
	if len(p.TopNotes) != 1 {
		t.Fatalf("expected duplicate/long notes to be removed, got %#v", p.TopNotes)
	}
}

func TestNormalizeProductPayloadRejectsBadInputs(t *testing.T) {
	cases := []struct {
		name string
		p    productPayload
	}{
		{
			name: "bad slug",
			p: productPayload{
				Name: "Fragancia", Slug: "Fragancia Nueva", ScentFamily: "Fresco",
				Variants: []variantForm{{SizeML: 50, PriceARSCents: 1, Stock: 1}},
				Images:   []imageForm{{URL: "https://example.com/a.jpg"}},
			},
		},
		{
			name: "negative stock",
			p: productPayload{
				Name: "Fragancia", Slug: "fragancia", ScentFamily: "Fresco",
				Variants: []variantForm{{SizeML: 50, PriceARSCents: 1, Stock: -1}},
				Images:   []imageForm{{URL: "https://example.com/a.jpg"}},
			},
		},
		{
			name: "missing image",
			p: productPayload{
				Name: "Fragancia", Slug: "fragancia", ScentFamily: "Fresco",
				Variants: []variantForm{{SizeML: 50, PriceARSCents: 1, Stock: 1}},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := normalizeProductPayload(&tc.p); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestNormalizeDiscountPayload(t *testing.T) {
	active := true
	req := discountWriteRequest{
		Code:          " verano10 ",
		DiscountType:  "percent",
		DiscountValue: 10,
		MinOrderCents: 1000,
		Active:        &active,
	}
	if err := normalizeDiscountPayload(&req, true); err != nil {
		t.Fatalf("expected valid discount: %v", err)
	}
	if req.Code != "VERANO10" {
		t.Fatalf("expected code normalization, got %q", req.Code)
	}

	expired := time.Now().Add(-time.Hour)
	req.ExpiresAt = &expired
	if err := normalizeDiscountPayload(&req, true); err == nil {
		t.Fatal("expected expired discount to fail")
	}
}

func TestNormalizeHomepagePayloadDefaultsHeroRotation(t *testing.T) {
	req := homepageRequest{
		HeroHeading:       " Perfumería argentina ",
		HeroImageURL:      "https://example.com/hero.jpg",
		HeroCTALabel:      "Catálogo",
		HeroCTAURL:        "/fragrances",
		EditorialImageURL: "https://example.com/editorial.jpg",
	}

	if err := normalizeHomepagePayload(&req); err != nil {
		t.Fatalf("expected valid homepage: %v", err)
	}
	if req.HeroHeading != "Perfumería argentina" {
		t.Fatalf("expected heading trim, got %q", req.HeroHeading)
	}
	if req.HeroImageMode != defaultHeroImageMode {
		t.Fatalf("expected default hero image mode, got %q", req.HeroImageMode)
	}
	if req.HeroRotationIntervalMS != defaultHeroRotationIntervalMS {
		t.Fatalf("expected default hero interval, got %d", req.HeroRotationIntervalMS)
	}
}

func TestNormalizeHomepagePayloadRejectsHeroRotationInputs(t *testing.T) {
	cases := []homepageRequest{
		{
			HeroHeading:            "Hero",
			HeroImageURL:           "https://example.com/hero.jpg",
			HeroImageMode:          "unknown",
			HeroRotationIntervalMS: defaultHeroRotationIntervalMS,
			EditorialImageURL:      "https://example.com/editorial.jpg",
		},
		{
			HeroHeading:            "Hero",
			HeroImageURL:           "https://example.com/hero.jpg",
			HeroImageMode:          "product_covers",
			HeroRotationIntervalMS: 500,
			EditorialImageURL:      "https://example.com/editorial.jpg",
		},
	}

	for _, tc := range cases {
		if err := normalizeHomepagePayload(&tc); err == nil {
			t.Fatalf("expected homepage validation error for %#v", tc)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	if err := validatePassword("abc1234567", "admin"); err != nil {
		t.Fatalf("expected valid password: %v", err)
	}
	for _, password := range []string{"short1", "sololetraslargas", "1234567890", "admin123456"} {
		if err := validatePassword(password, "admin"); err == nil {
			t.Fatalf("expected %q to fail", password)
		}
	}
}

func TestNormalizeShippingZonePayload(t *testing.T) {
	req := shippingZonePayload{
		ZoneName:         " Interior ",
		ProvinceCodes:    []string{" ar, ba ", "BA"},
		BaseCostCents:    420000,
		EstimatedDaysMin: 3,
		EstimatedDaysMax: 7,
		Active:           true,
	}
	if err := normalizeShippingZonePayload(&req); err != nil {
		t.Fatalf("expected valid zone: %v", err)
	}
	if req.ZoneName != "Interior" {
		t.Fatalf("expected trimmed name, got %q", req.ZoneName)
	}
	if len(req.ProvinceCodes) != 2 || req.ProvinceCodes[0] != "AR" || req.ProvinceCodes[1] != "BA" {
		t.Fatalf("unexpected province codes: %#v", req.ProvinceCodes)
	}
}
