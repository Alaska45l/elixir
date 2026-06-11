package config

import (
	"reflect"
	"testing"
)

func TestAllowedOriginsIncludesFrontendURL(t *testing.T) {
	got := allowedOrigins("https://old-front.onrender.com", "https://new-front.onrender.com")
	want := []string{"https://old-front.onrender.com", "https://new-front.onrender.com"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("allowedOrigins() = %#v, want %#v", got, want)
	}
}

func TestAllowedOriginsNormalizesAndDeduplicates(t *testing.T) {
	got := allowedOrigins(" https://front.onrender.com/ , https://admin.example ", "https://front.onrender.com")
	want := []string{"https://front.onrender.com", "https://admin.example"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("allowedOrigins() = %#v, want %#v", got, want)
	}
}
