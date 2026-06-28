# SpotSync

SpotSync is a full-stack smart parking and EV charging reservation system for airports and malls.

## Live URLs

- Backend API: https://spotsync-api-8acf.onrender.com
- Frontend: https://project-a7ord.vercel.app

## Projects

- `spotsync-api`: Go, Echo, GORM, PostgreSQL backend.
- `spotsync-web`: React, TypeScript, Vite frontend.

## Features

- Register and login with JWT authentication.
- Driver and admin role support.
- Public parking zone browsing with dynamic availability.
- Transaction-safe parking reservations using PostgreSQL row locks.
- Driver reservation list and cancellation.
- Admin dashboard for zones and reservations.

## Architecture

Backend requests follow strict clean architecture:

```text
HTTP request
  -> handler: bind, validate, authorize, respond
  -> service: business rules, JWT, bcrypt, ownership checks
  -> repository: GORM database operations, preloads, transactions
  -> PostgreSQL
```

## Backend Setup

```bash
cd spotsync-api
cp .env.example .env
go mod tidy
go run .
```

For a quick local PostgreSQL database, run this from the repository root:

```bash
docker compose up -d
```

Then use this local database URL in `spotsync-api/.env`:

```env
DATABASE_URL=postgresql://spotsync:spotsync_pass@localhost:5432/spotsync?sslmode=disable
```

Required environment variables:

```env
DATABASE_URL=postgresql://user:password@host/spotsync?sslmode=require
JWT_SECRET=replace-with-a-long-random-secret-at-least-32-chars
PORT=8080
FRONTEND_URL=http://localhost:5173
```

## Frontend Setup

```bash
cd spotsync-web
cp .env.example .env
npm install
npm run dev
```

Frontend environment:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## API Endpoints

| Method | Endpoint | Access |
| --- | --- | --- |
| POST | `/api/v1/auth/register` | Public |
| POST | `/api/v1/auth/login` | Public |
| GET | `/api/v1/zones` | Public |
| GET | `/api/v1/zones/:id` | Public |
| POST | `/api/v1/zones` | Admin |
| PUT | `/api/v1/zones/:id` | Admin |
| DELETE | `/api/v1/zones/:id` | Admin |
| POST | `/api/v1/reservations` | Authenticated |
| GET | `/api/v1/reservations/my-reservations` | Authenticated |
| DELETE | `/api/v1/reservations/:id` | Authenticated |
| GET | `/api/v1/reservations` | Admin |

## Verification

```bash
cd spotsync-api && go test ./...
cd spotsync-web && npm run build
```

## Deployment

- Backend: Render with `DATABASE_URL`, `JWT_SECRET`, and `FRONTEND_URL`.
- Database: NeonDB, Supabase, or Aiven PostgreSQL.
- Frontend: Vercel with `VITE_API_URL` pointing to the deployed backend.

Included deployment helpers:

- `spotsync-api/render.yaml`
- `spotsync-web/vercel.json`
