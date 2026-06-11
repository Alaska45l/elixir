package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elixir/backend/internal/admin"
	"elixir/backend/internal/config"
	"elixir/backend/internal/contact"
	"elixir/backend/internal/db"
	"elixir/backend/internal/discount"
	"elixir/backend/internal/health"
	"elixir/backend/internal/media"
	"elixir/backend/internal/middleware"
	"elixir/backend/internal/orders"
	"elixir/backend/internal/payments"
	"elixir/backend/internal/products"
	"elixir/backend/internal/shipping"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := cfg.Validate(); err != nil {
		slog.Error("config validation failed", "error", err)
		os.Exit(1)
	}
	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if pool != nil {
		defer pool.Close()
	}

	loginLimiter := middleware.NewLoginLimiter(5, 10*time.Minute)
	apiLimiter := middleware.NewAPILimiter(60, time.Minute)
	webhookLimiter := middleware.NewAPILimiter(120, time.Minute)
	secureCookies := admin.IsSecureEnv(cfg.AppEnv)
	sessions := admin.SessionManager{Secret: cfg.SessionSecret, Duration: cfg.SessionDuration, Secure: secureCookies, SameSite: admin.ResolveSameSite(cfg.SessionSameSite, secureCookies)}
	discountSvc := discount.Service{Repo: discount.Repository{Pool: pool}}
	orderSvc := orders.Service{Repo: orders.DBRepository{Pool: pool}, Discount: discountSvc}
	storageSvc, err := media.NewStorageService(cfg.R2AccountID, cfg.R2AccessKey, cfg.R2SecretKey, cfg.R2BucketName, cfg.R2PublicURL)
	if err != nil {
		log.Fatal("media storage:", err)
	}

	productHandler := products.Handler{Service: products.Service{Repo: products.Repository{Pool: pool}}}
	orderHandler := orders.Handler{Service: orderSvc}
	paymentHandler := payments.Handler{
		Service: payments.Service{
			Repo: payments.DBRepository{Pool: pool},
			MP:   payments.MercadoPagoClient{AccessToken: cfg.MPAccessToken, BackendURL: cfg.BackendURL, FrontendURL: cfg.FrontendURL},
		},
		WebhookSecret: cfg.MPWebhookSecret,
	}
	adminHandler := admin.Handler{Pool: pool, Sessions: sessions, Auth: admin.AuthService{Pool: pool, Sessions: sessions, Limiter: loginLimiter}, Media: storageSvc}
	contactHandler := contact.Handler{Service: contact.Service{Pool: pool}}
	shippingProviders := []shipping.ShippingProvider{
		shipping.LocalPickupProvider{},
		shipping.CorreoArgentinoProvider{APIKey: cfg.CorreoArgAPIKey, ClientID: cfg.CorreoArgClientID, OriginPostalCode: cfg.OriginPostalCode},
		&shipping.AndreaniProvider{User: cfg.AndreaniUser, Password: cfg.AndreaniPassword, ClientID: cfg.AndreaniClientID, OriginPostalCode: cfg.OriginPostalCode},
	}
	shippingHandler := shipping.Handler{Service: shipping.Service{Repo: shipping.Repository{Pool: pool}, Providers: shippingProviders}}

	r := chi.NewRouter()
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.Compress)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestTimeout(15 * time.Second))
	r.Use(middleware.RequestLogger)
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	csrf := middleware.CSRF(cfg.AllowedOrigins, cfg.BackendURL)

	r.Group(func(ur chi.Router) {
		ur.Use(middleware.NoStore, csrf, middleware.AdminSession(sessions))
		ur.Post("/api/admin/upload", adminHandler.UploadImage)
	})

	r.Group(func(rr chi.Router) {
		rr.Use(middleware.BodyLimit(2 << 20))
		rr.Use(csrf)

		rr.Mount("/", health.Handler{Pool: pool, Started: time.Now(), Version: "1.0.0"}.Routes())
		rr.Get("/api/homepage", adminHandler.PublicHomepage)
		rr.Get("/api/settings", adminHandler.PublicSettings)
		rr.With(apiLimiter.Throttle).Get("/api/products/search", productHandler.Search)
		rr.Mount("/api/products", productHandler.Routes())
		rr.With(apiLimiter.Throttle).Post("/api/cart/validate", orderHandler.ValidateCart)
		rr.With(apiLimiter.Throttle).Post("/api/discount/validate", discount.Handler{Service: discountSvc}.Validate)
		rr.With(apiLimiter.Throttle).Post("/api/orders", orderHandler.Create)
		rr.Get("/api/orders/{external_reference}", orderHandler.Status)
		rr.With(apiLimiter.Throttle).Post("/api/checkout/mercadopago/preference", paymentHandler.CreatePreference)
		rr.With(webhookLimiter.Throttle).Post("/api/payments/mercadopago/webhook", paymentHandler.Webhook)
		rr.Get("/api/shipping/zones", shippingHandler.Zones)
		rr.With(apiLimiter.Throttle).Post("/api/shipping/quote", shippingHandler.Quote)
		rr.With(apiLimiter.Throttle).Post("/api/contact", contactHandler.Message)
		rr.With(apiLimiter.Throttle).Post("/api/contact/abandoned-cart", contactHandler.AbandonedCart)

		rr.With(middleware.NoStore, loginLimiter.Middleware).Post("/api/admin/login", adminHandler.Login)
		rr.With(middleware.NoStore).Post("/api/admin/logout", adminHandler.Logout)
		rr.Group(func(ar chi.Router) {
			ar.Use(middleware.NoStore)
			ar.Use(middleware.AdminSession(sessions))
			ar.Get("/api/admin/me", adminHandler.Me)
			ar.Post("/api/admin/password", adminHandler.ChangePassword)
			ar.Get("/api/admin/users", adminHandler.AdminUsers)
			ar.Post("/api/admin/users", adminHandler.AdminUsers)
			ar.Post("/api/admin/users/{username}/password", adminHandler.AdminUserByUsername)
			ar.Delete("/api/admin/users/{username}", adminHandler.AdminUserByUsername)
			ar.Get("/api/admin/metrics", adminHandler.Metrics)
			ar.Get("/api/admin/products", adminHandler.AdminProducts)
			ar.Post("/api/admin/products", adminHandler.SaveProduct)
			ar.Get("/api/admin/products/{id}", adminHandler.AdminProductByID)
			ar.Put("/api/admin/products/{id}", adminHandler.SaveProduct)
			ar.Put("/api/admin/products/{id}/active", adminHandler.UpdateProductActive)
			ar.Delete("/api/admin/products/{id}", adminHandler.DeleteProduct)
			ar.Post("/api/admin/products/import", adminHandler.ImportProducts)
			ar.Get("/api/admin/orders", adminHandler.Orders)
			ar.Get("/api/admin/orders/export", adminHandler.ExportOrders)
			ar.Put("/api/admin/orders/{id}", adminHandler.UpdateOrder)
			ar.Get("/api/admin/discounts", adminHandler.Discounts)
			ar.Post("/api/admin/discounts", adminHandler.Discounts)
			ar.Put("/api/admin/discounts/{id}", adminHandler.DiscountByID)
			ar.Delete("/api/admin/discounts/{id}", adminHandler.DiscountByID)
			ar.Get("/api/admin/homepage", adminHandler.Homepage)
			ar.Put("/api/admin/homepage", adminHandler.Homepage)
			ar.Get("/api/admin/settings", adminHandler.Settings)
			ar.Put("/api/admin/settings", adminHandler.Settings)
			ar.Get("/api/admin/shipping/zones", adminHandler.AdminShippingZones)
			ar.Post("/api/admin/shipping/zones", adminHandler.AdminShippingZones)
			ar.Put("/api/admin/shipping/zones/{id}", adminHandler.AdminShippingZoneByID)
			ar.Delete("/api/admin/shipping/zones/{id}", adminHandler.AdminShippingZoneByID)
			ar.Get("/api/admin/contact", adminHandler.Contact)
			ar.Put("/api/admin/contact/{id}/read", adminHandler.MarkContactRead)
			ar.Get("/api/admin/low-stock", adminHandler.LowStock)
		})
	})

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r, ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	slog.Info("api listening", "port", cfg.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	slog.Info("api shutting down")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
