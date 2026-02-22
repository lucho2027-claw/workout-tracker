# Workout Tracker (Monorepo)

Fullstack workout tracking app.

## Stack
- Backend: Go stdlib `net/http` + PostgreSQL
- Frontend: Next.js (App Router) + TypeScript

## Monorepo layout
- `apps/api` Go API
- `apps/web` Next.js frontend

## Quick start
```bash
docker compose up -d db
cd apps/api && cp .env.example .env && go run ./cmd/server
cd ../web && cp .env.local.example .env.local && npm install && npm run dev
```

## API base URL
`http://localhost:8080`

## Web URL
`http://localhost:3000`
