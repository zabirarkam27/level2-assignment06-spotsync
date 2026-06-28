# SpotSync

SpotSync is a clean-architecture REST API for smart parking and EV charging reservations at busy airports and malls.

## Live URL

- Backend API: https://spotsync-api-8acf.onrender.com

## Features

- User registration and login with JWT authentication.
- Password hashing with bcrypt.
- Driver and admin role-based authorization.
- Public parking zone listing with dynamic available spot calculation.
- Admin parking zone create, update, delete, and all-reservation access.
- Authenticated reservation creation, own-reservation listing, and cancellation.
- Transaction-safe reservation creation with PostgreSQL row-level locking to prevent overbooking.
- Centralized API response and error formatting.

## Tech Stack

- Go 1.22+
- Echo
- GORM
- PostgreSQL
- go-playground/validator
- golang-jwt/jwt/v5
- bcrypt

## Architecture

The project follows strict clean architecture. Each layer has a separate responsibility:

```text
HTTP Request
  -> Handler: bind request DTOs, validate input, read JWT context, return JSON responses
  -> Service: business rules, password hashing, JWT creation, ownership and capacity checks
  -> Repository: GORM queries, transactions, row locks, preloads
  -> PostgreSQL
```

Project folders:

```text
config/       database connection
dto/          request and response DTOs
handler/      HTTP handlers
middleware/   JWT and role middleware
models/       GORM models
repository/   database access layer
service/      business logic layer
```

Dependencies are wired manually in `main.go`:

```text
Repository -> Service -> Handler -> Routes
```

## Environment Variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgresql://user:password@host/database?sslmode=require
JWT_SECRET=replace-with-a-long-random-secret-at-least-32-chars
PORT=8080
CORS_ALLOWED_ORIGINS=*
```

For local Docker PostgreSQL, run:

```bash
docker compose up -d
```

Then use:

```env
DATABASE_URL=postgresql://spotsync:spotsync_pass@localhost:55432/spotsync?sslmode=disable
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

## Request Examples

Register:

```json
{
  "name": "John Doe",
  "email": "john.doe@spotsync.com",
  "password": "securePassword123",
  "role": "driver"
}
```

Login:

```json
{
  "email": "john.doe@spotsync.com",
  "password": "securePassword123"
}
```

Create zone:

```json
{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 20,
  "price_per_hour": 5.5
}
```

Create reservation:

```json
{
  "zone_id": 5,
  "license_plate": "ABC-1234"
}
```

## Response Format

Success response:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {}
}
```

Error response:

```json
{
  "success": false,
  "message": "Error description",
  "errors": "Error details"
}
```

## Verification

```bash
go test ./...
```

## Deployment

The backend is deployed on Render and uses NeonDB PostgreSQL. Required deployment variables:

```env
DATABASE_URL=postgresql://user:password@host/database?sslmode=require
JWT_SECRET=replace-with-a-long-random-secret-at-least-32-chars
CORS_ALLOWED_ORIGINS=*
```

`render.yaml` is included for the backend service configuration.
