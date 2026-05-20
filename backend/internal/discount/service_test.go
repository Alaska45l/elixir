package discount

import (
	"testing"
	"time"
)

func TestValidateCodeExpired(t *testing.T) {
	expired := time.Now().Add(-time.Hour)
	res := ValidateCode(Code{Code: "OLD", Active: true, DiscountType: "percent", DiscountValue: 10, ExpiresAt: &expired}, 10000, time.Now())
	if res.Valid {
		t.Fatal("expected expired code to be invalid")
	}
}

func TestValidateCodeMaxUsesReached(t *testing.T) {
	max := 1
	res := ValidateCode(Code{Code: "MAX", Active: true, DiscountType: "percent", DiscountValue: 10, MaxUses: &max, Uses: 1}, 10000, time.Now())
	if res.Valid {
		t.Fatal("expected max uses reached")
	}
}

func TestValidateCodeMinOrderNotMet(t *testing.T) {
	res := ValidateCode(Code{Code: "MIN", Active: true, DiscountType: "fixed", DiscountValue: 1000, MinOrderCents: 50000}, 10000, time.Now())
	if res.Valid {
		t.Fatal("expected min order rejection")
	}
}
