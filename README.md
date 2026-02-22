# Workout Tracker (Monorepo)

Fullstack workout tracking app.

## Stack
- Backend: Go stdlib `net/http` + PostgreSQL
- Frontend: Next.js (App Router) + TypeScript + Tailwind v4

## Monorepo layout
- `apps/api` Go API
- `apps/web` Next.js frontend

## Local development
1) Start PostgreSQL (Docker):
```bash
docker compose up -d db
```

2) Run API:
```bash
cd apps/api
cp .env.example .env
# set a secure JWT_SECRET before production
go run ./cmd/server
```

3) Run Web:
```bash
cd apps/web
cp .env.local.example .env.local
npm install
npm run dev
```

## Production (containerized)
```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Before production:
- Set a strong `JWT_SECRET`
- Set `ALLOWED_ORIGIN` to your frontend URL/domain
- Put services behind HTTPS reverse proxy (Nginx/Caddy)
- Do not expose PostgreSQL publicly

## API base URL
`http://localhost:8080`

## Web URL
`http://localhost:3000`
