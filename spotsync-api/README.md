# SpotSync API

SpotSync is a clean-architecture Go API for smart parking and EV charging spot reservations.

Live API: https://spotsync-api-8acf.onrender.com

## Features

- Driver/admin authentication with JWT and bcrypt.
- Public parking zone availability.
- Admin zone CRUD and all-reservation view.
- Driver reservation creation, listing, and cancellation.
- Transactional reservation creation with row-level locking to prevent overbooking.

## Tech Stack

- Go 1.22+
- Echo
- GORM
- PostgreSQL
- go-playground/validator
- golang-jwt/jwt/v5
- bcrypt

## Architecture

Requests flow through `handler -> service -> repository -> database`.

- `dto`: request and response contracts.
- `handler`: HTTP binding, validation, auth context, response formatting.
- `service`: business rules, password hashing, JWT generation, ownership checks.
- `repository`: all GORM queries, transactions, row locks, and preloads.
- `models`: GORM database tables.

## Environment

Copy `.env.example` to `.env`.

```env
DATABASE_URL=postgresql://user:password@host/spotsync?sslmode=require
JWT_SECRET=replace-with-a-long-random-secret-at-least-32-chars
PORT=8080
FRONTEND_URL=http://localhost:5173
```

## Run Locally

```bash
go mod tidy
go run .
```

Health check:

```bash
GET http://localhost:8080/api/v1/health
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

## Concurrency Note

`repository/reservation_repository.go` creates reservations inside a GORM transaction. It locks the selected parking zone row using `clause.Locking{Strength: "UPDATE"}`, counts active reservations, checks capacity, and creates the reservation atomically.
