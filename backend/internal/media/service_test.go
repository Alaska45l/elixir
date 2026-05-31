package media

import "testing"

func TestNewStorageServiceNoConfigIsOptional(t *testing.T) {
	service, err := NewStorageService("", "", "", "", "")
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}
	if service != nil {
		t.Fatal("expected nil service when R2 is not configured")
	}
}

func TestNewStorageServiceRejectsPartialConfig(t *testing.T) {
	service, err := NewStorageService("account", "access", "", "bucket", "https://cdn.example.com")
	if err == nil {
		t.Fatal("expected partial config error")
	}
	if service != nil {
		t.Fatal("expected nil service on partial config")
	}
}

func TestSanitizeFolder(t *testing.T) {
	got := sanitizeFolder("../Products/Hero Images")
	if got != "products/heroimages" {
		t.Fatalf("expected sanitized folder, got %q", got)
	}
}
