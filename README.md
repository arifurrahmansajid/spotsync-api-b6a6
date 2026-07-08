# SpotSync – Smart Parking & EV Charging Reservation

A clean-architecture REST API built with **Go**, **Echo**, **GORM**, and **PostgreSQL**.

## Architecture

```
spotsync/
├── config/          # Database connection & migration
├── dto/             # Request/Response shapes (never exposes GORM models)
├── handler/         # HTTP layer – binds DTOs, calls Service, returns JSON
├── middleware/       # JWT auth & role enforcement
├── models/          # GORM structs (database tables)
├── repository/      # All GORM operations (CRUD + transactions)
├── service/         # Business logic (hashing, JWT signing, capacity check)
├── validator/       # Echo-integrated go-playground/validator
└── main.go          # Dependency injection entrypoint
```

## Setup

### 1. Prerequisites
- Go 1.22+
- PostgreSQL (or [NeonDB](https://neon.tech) / [Supabase](https://supabase.com))

### 2. Configure Environment

```bash
cp .env.example .env
# Edit .env with your database credentials and JWT secret
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run

```bash
# With godotenv (install once: go install github.com/joho/godotenv/cmd/godotenv@latest)
godotenv -f .env go run main.go

# Or export variables manually
$env:DATABASE_URL="postgres://..."
$env:JWT_SECRET="your-secret"
go run main.go
```

The server starts on **http://localhost:8080** by default.

## API Endpoints

### Authentication (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login and receive JWT |

### Parking Zones
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| GET | `/api/v1/zones` | Public | List all zones with available spots |
| GET | `/api/v1/zones/:id` | Public | Get single zone |
| POST | `/api/v1/zones` | Admin | Create zone |
| PUT | `/api/v1/zones/:id` | Admin | Update zone |
| DELETE | `/api/v1/zones/:id` | Admin | Delete zone |

### Reservations
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| POST | `/api/v1/reservations` | Authenticated | Reserve a spot |
| GET | `/api/v1/reservations/my-reservations` | Authenticated | My reservations |
| DELETE | `/api/v1/reservations/:id` | Authenticated | Cancel reservation |
| GET | `/api/v1/reservations` | Admin | All reservations |

## Key Implementation Details

### Concurrency Safety (EV Spot Bottleneck)
Reservations use **SELECT FOR UPDATE** row-level locking inside a GORM transaction:
```
db.Transaction(func(tx) {
  tx.Clauses(clause.Locking{Strength:"UPDATE"}).First(&zone, id)
  // count active reservations
  // reject if full (409 Conflict)
  // else create reservation
})
```
This prevents race conditions where two drivers could both book the last available spot simultaneously.

### JWT Payload
```json
{ "id": 1, "role": "driver", "exp": 1234567890 }
```

### Password Security
- bcrypt cost factor: **10**
- Passwords **never** appear in API responses or logs
