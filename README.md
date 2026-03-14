# Event Budaya Ticketing - BCC

A ticketing platform for cultural events. This application provides APIs for event management, categories, authentication, and ticketing system.

## 📋 Prerequisites

Before starting, make sure you have installed:

- **Go** v1.25.0 or higher ([Download](https://golang.org/dl))
- **PostgreSQL** v12 or higher ([Download](https://www.postgresql.org/download))
- **Git**
- **Make** (optional, for command shortcuts)

Verify the installation:
```bash
go version
psql --version
git --version
```

## 🚀 Installation & Setup

### 1. Clone Repository

```bash
git clone https://github.com/rafirh/event-budaya-ticketing-bcc.git
cd event-budaya-ticketing-bcc
```

### 2. Download Dependencies

```bash
go mod download
go mod tidy
```

### 3. Setup Environment Variables

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

Edit `.env` according to your local configuration:

```env
# App Configuration
APP_NAME=event-budaya-ticketing
APP_ENV=development
APP_PORT=3000

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=event_budaya_db

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this
JWT_EXPIRY_HOURS=24

# S3 Configuration (optional)
S3_KEY=your_aws_access_key
S3_SECRET=your_aws_secret_key
S3_BUCKET=your_bucket_name
S3_REGION=ap-southeast-1
S3_PUBLIC_BASE=https://your-cloudfront-url.com
```

### 4. Setup Database

#### Create PostgreSQL database:

```bash
createdb event_budaya_db
```

Or using psql:
```bash
psql -U postgres
postgres=# CREATE DATABASE event_budaya_db;
postgres=# \q
```

#### Run migrations:

```bash
go run cmd/main.go -migrate
```

#### Seed data (optional):

```bash
go run cmd/main.go -seed
```

**Or fresh setup (drop & recreate):**

```bash
go run cmd/main.go -fresh
```

### 5. Build & Run

**Development Mode:**
```bash
go run cmd/main.go
```

**Build Binary:**
```bash
go build -o bin/app cmd/main.go
./bin/app
```

Server will run at `http://localhost:3000`
API documentation available at `https://2ikhh28mj3.apidog.io/`

**Last Updated:** March 14, 2026  
**Version:** 1.0.0
