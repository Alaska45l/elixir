# ELIXIR Exclusive

Luxury Argentine perfume e-commerce store with SvelteKit SSR frontend, Go API backend, PostgreSQL migrations for Neon, MercadoPago checkout, admin panel, seed data and deployment notes.

## Commands to run locally

```bash
# Backend
cd backend
cp .env.example .env
# edit DATABASE_URL, SESSION_SECRET and MP_ACCESS_TOKEN
go mod tidy
go run ./cmd/migrate
go run ./cmd/seed
go run ./cmd/api
```

```bash
# Frontend
cd frontend
cp .env.example .env
npm install
npm run dev
```

```bash
# Verification
cd frontend && npm run check && npm run build
cd backend && go test ./...
```

## Neon

1. Create a Neon project.
2. Copy the pooled PostgreSQL connection string.
3. Set `DATABASE_URL=postgresql://...?...sslmode=require`.
4. Run `go run ./cmd/migrate`.
5. Run `go run ./cmd/seed`.

## Render backend

1. Create a Go web service from this repository.
2. Build command: `go build -o app ./cmd/api/main.go`.
3. Start command: `./app`.
4. Health check path: `/api/health`.
5. Add all backend environment variables from `backend/.env.example`.

## Render frontend

1. Create a Node web service using `frontend` as the root directory.
2. Build command: `npm run build`.
3. Start command: `node build`.
4. Set `PUBLIC_API_URL` to the backend Render URL.
5. Set `PUBLIC_SITE_URL` to the frontend production URL.

## MercadoPago

1. Create an app in the MercadoPago Argentina developer portal.
2. Set webhook URL to `https://your-backend.render.com/api/payments/mercadopago/webhook`.
3. Subscribe to `payment` events.
4. Store the access token only in backend `MP_ACCESS_TOKEN`.

## Manual operator steps

- Change the seeded admin password immediately.
- Replace demo product images with licensed final product photography.
- Configure the real MercadoPago production token before accepting payments.
- Configure a real WhatsApp number and production domain.

## Known limitations

- SMTP variables are documented but email sending is not active; contact messages are stored in the database and shown in admin.
- CSV import endpoint returns an acknowledgement stub; the admin product form supports direct product creation.
- This workspace does not currently expose a Go runtime, so backend tests could not be executed here. Frontend `svelte-check` and build were executed successfully.
