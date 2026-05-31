# ELIXIR Exclusive Backend

Go API for catalog, cart validation, orders, MercadoPago checkout, admin sessions, shipping and contact messages.

## Local

```bash
cp .env.example .env
export $(grep -v '^#' .env | xargs)
go run ./cmd/migrate
go run ./cmd/seed
go run ./cmd/api
```

The seed creates `admin` / `elixir2024`. Change this password immediately after first login.

## Smoke Tests

```bash
curl http://localhost:8080/api/health
curl http://localhost:8080/api/products
curl http://localhost:8080/api/products/nocturno-oud
```

## Render

- Build command: `go build -o app ./cmd/api/main.go`
- Start command: `./app`
- Health check path: `/api/health`
- Set all variables from `.env.example` in the Render dashboard.
- Image uploads require Cloudflare R2 variables. `R2_PUBLIC_URL` should be the public bucket/custom-domain base URL, without a trailing slash.
- WebP encoding uses CGO, so the build environment must have a C compiler available.
