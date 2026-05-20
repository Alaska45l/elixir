package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"elixir/backend/internal/admin"
	"elixir/backend/internal/config"
	"elixir/backend/internal/contact"
	"elixir/backend/internal/db"
	"elixir/backend/internal/discount"
	"elixir/backend/internal/health"
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
	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if pool != nil {
		defer pool.Close()
	}

	loginLimiter := middleware.NewLoginLimiter(5, 10*time.Minute)
	sessions := admin.SessionManager{Secret: cfg.SessionSecret, Duration: cfg.SessionDuration, Secure: admin.IsSecureEnv(cfg.AppEnv)}
	discountSvc := discount.Service{Repo: discount.Repository{Pool: pool}}
	orderSvc := orders.Service{Repo: orders.DBRepository{Pool: pool}, Discount: discountSvc}

	productHandler := products.Handler{Service: products.Service{Repo: products.Repository{Pool: pool}}}
	orderHandler := orders.Handler{Service: orderSvc}
	paymentHandler := payments.Handler{Service: payments.Service{
		Repo: payments.DBRepository{Pool: pool},
		MP:   payments.MercadoPagoClient{AccessToken: cfg.MPAccessToken, BackendURL: cfg.BackendURL, FrontendURL: cfg.FrontendURL},
	}}
	adminHandler := admin.Handler{Pool: pool, Sessions: sessions, Auth: admin.AuthService{Pool: pool, Sessions: sessions, Limiter: loginLimiter}}
	contactHandler := contact.Handler{Service: contact.Service{Pool: pool}}
	shippingHandler := shipping.Handler{Service: shipping.Service{Repo: shipping.Repository{Pool: pool}}}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(middleware.BodyLimit(2 << 20))

	r.Mount("/", health.Handler{Pool: pool, Started: time.Now(), Version: "1.0.0"}.Routes())
	r.Mount("/api/products", productHandler.Routes())
	r.Post("/api/cart/validate", orderHandler.ValidateCart)
	r.Post("/api/discount/validate", discount.Handler{Service: discountSvc}.Validate)
	r.Post("/api/orders", orderHandler.Create)
	r.Get("/api/orders/{external_reference}", orderHandler.Status)
	r.Post("/api/checkout/mercadopago/preference", paymentHandler.CreatePreference)
	r.Post("/api/payments/mercadopago/webhook", paymentHandler.Webhook)
	r.Get("/api/shipping/zones", shippingHandler.Zones)
	r.Post("/api/contact", contactHandler.Message)
	r.Post("/api/contact/abandoned-cart", contactHandler.AbandonedCart)

	r.With(loginLimiter.Middleware).Post("/api/admin/login", adminHandler.Login)
	r.Post("/api/admin/logout", adminHandler.Logout)
	r.Group(func(ar chi.Router) {
		ar.Use(middleware.AdminSession(sessions))
		ar.Get("/api/admin/me", adminHandler.Me)
		ar.Get("/api/admin/metrics", adminHandler.Metrics)
		ar.Get("/api/admin/products", adminHandler.AdminProducts)
		ar.Post("/api/admin/products", adminHandler.SaveProduct)
		ar.Put("/api/admin/products/{id}", adminHandler.SaveProduct)
		ar.Delete("/api/admin/products/{id}", adminHandler.DeleteProduct)
		ar.Post("/api/admin/products/import", adminHandler.ImportProducts)
		ar.Get("/api/admin/orders", adminHandler.Orders)
		ar.Put("/api/admin/orders/{id}", adminHandler.UpdateOrder)
		ar.Get("/api/admin/discounts", adminHandler.Discounts)
		ar.Post("/api/admin/discounts", adminHandler.Discounts)
		ar.Put("/api/admin/discounts/{id}", adminHandler.DiscountByID)
		ar.Delete("/api/admin/discounts/{id}", adminHandler.DiscountByID)
		ar.Get("/api/admin/homepage", adminHandler.Homepage)
		ar.Put("/api/admin/homepage", adminHandler.Homepage)
		ar.Get("/api/admin/contact", adminHandler.Contact)
		ar.Put("/api/admin/contact/{id}/read", adminHandler.MarkContactRead)
		ar.Get("/api/admin/low-stock", adminHandler.LowStock)
	})

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r, ReadHeaderTimeout: 5 * time.Second}
	slog.Info("api listening", "port", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
