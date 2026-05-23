package admin

import (
	"net/http"
	"testing"
)

func TestResolveSameSite(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		secure bool
		want   http.SameSite
	}{
		{name: "default strict", value: "", secure: true, want: http.SameSiteStrictMode},
		{name: "strict", value: "strict", secure: true, want: http.SameSiteStrictMode},
		{name: "lax", value: "lax", secure: true, want: http.SameSiteLaxMode},
		{name: "none secure", value: "none", secure: true, want: http.SameSiteNoneMode},
		{name: "none insecure falls back", value: "none", secure: false, want: http.SameSiteLaxMode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ResolveSameSite(tt.value, tt.secure); got != tt.want {
				t.Fatalf("ResolveSameSite(%q, %v) = %v, want %v", tt.value, tt.secure, got, tt.want)
			}
		})
	}
}
