# Workout Tracker (Monorepo)

Fullstack workout tracking app.

## Stack
- Backend: Go stdlib `net/http` + PostgreSQL
- Frontend: Next.js (App Router) + TypeScript + Tailwind v4
- SQL bridge: `sqlc` (configured in `apps/api/sqlc.yaml`)

## Monorepo layout
- `apps/api` Go API
- `apps/web` Next.js frontend

---

## Why login may fail locally
If you only try **Login** without creating a user first, auth will fail.

Use this order:
1. Start DB
2. Run migrations
3. Start API
4. Start Web
5. Open **Register** page and create account
6. Then login

---

## Local development (step-by-step)

### 1) Start PostgreSQL
```bash
docker compose up -d db
```

### 2) API env
```bash
cd apps/api
cp .env.example .env
# Set a secure JWT secret (even locally)
# JWT_SECRET=some-long-random-string
```

### 3) Run migrations
```bash
cd apps/api
go run ./cmd/migrate
```

### 4) Start API
```bash
cd apps/api
go run ./cmd/server
```

### 5) Start Web
```bash
cd apps/web
cp .env.local.example .env.local
npm install
npm run dev
```

### 6) Use the app
- Web: `http://localhost:3000`
- Register: `http://localhost:3000/register`
- Login: `http://localhost:3000/login`

---

## sqlc usage
`sqlc` is configured and query files are in:
- `apps/api/queries/*.sql`
- config: `apps/api/sqlc.yaml`

Generate Go code:
```bash
cd apps/api
sqlc generate
```

> Note: current handlers still include direct SQL calls. Next refactor step is to switch handlers to generated `sqlc` query methods completely.

---

## Production (containerized)
```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Before production:
- Set a strong `JWT_SECRET`
- Set `ALLOWED_ORIGIN` to your frontend URL/domain
- Put services behind HTTPS reverse proxy (Nginx/Caddy)
- Do not expose PostgreSQL publicly

