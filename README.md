# Event Budaya Ticketing - BCC

A ticketing platform for cultural events managed by Balai Budaya Citra (BCC). This application provides APIs for event management, categories, authentication, and ticketing system.

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

## 📚 API Endpoints

### Base URL
```
http://localhost:3000/api
```

### Health Check
```
GET /health
```

### Authentication
```
POST /auth/register    - Register new user
POST /auth/login       - Login & get JWT token
POST /auth/logout      - Logout (delete token)
```

**Login Request:**
```json
{
  "email": "user@gmail.com",
  "password": "user123"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGc...",
    "user": {
      "id": "uuid",
      "name": "User",
      "email": "user@gmail.com",
      "role": "user"
    }
  }
}
```

### Event Categories
```
GET /api/categories        - Get all categories
GET /api/categories/:id    - Get category by ID
```

**Category Response:**
```json
{
  "id": "uuid",
  "name": "Seni Pertunjukan",
  "logo": "https://event-budaya.iccn.or.id/seeders/categories/seni_pertunjukan.svg"
}
```

### Events
```
GET /api/events              - Get all events
GET /api/events/:slug        - Get event detail by slug
```

**Events Response:**
```json
{
  "status": "success",
  "message": "Events retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "title": "Seminar Budaya Nusantara 2026",
      "slug": "seminar-budaya-nusantara-2026",
      "description": "Educational session about Nusantara's cultural richness.",
      "promoter": {
        "id": "uuid",
        "name": "Promotor"
      },
      "category": {
        "id": "uuid",
        "name": "Seminar Budaya",
        "logo": "https://..."
      },
      "venue": "Gedung Kesenian Malang",
      "address": "Jl. Soekarno Hatta No. 1, Malang",
      "start_date": "2026-03-21T09:00:00Z",
      "end_date": "2026-03-21T12:00:00Z",
      "is_paid": true,
      "banner_url": null,
      "status": "published",
      "created_at": "2026-03-14T10:00:00Z"
    }
  ]
}
```

## 📁 Project Structure

```
event-budaya-ticketing-bcc/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   ├── config.go              # App & database configuration
│   └── database.go            # Database connection setup
├── db/
│   └── db.sql                 # Database schema
├── internal/
│   ├── dto/                   # Data Transfer Objects
│   │   ├── event_response.go
│   │   ├── user_request.go
│   │   └── user_response.go
│   ├── handler/               # HTTP request handlers
│   │   ├── auth_handler.go
│   │   ├── category_handler.go
│   │   └── event_handler.go
│   ├── middleware/            # Custom middlewares
│   │   ├── auth.go
│   │   └── logger.go
│   ├── model/                 # Database models
│   │   ├── event.go
│   │   ├── event_category.go
│   │   ├── personal_access_token.go
│   │   └── user.go
│   ├── repository/            # Data access layer
│   │   ├── interfaces.go
│   │   └── gorm/
│   │       ├── auth_repository.go
│   │       ├── event_category_repository.go
│   │       ├── event_repository.go
│   │       └── personal_access_token_repository.go
│   ├── router/                # Route definitions
│   │   ├── auth_route.go
│   │   ├── category_route.go
│   │   ├── event_route.go
│   │   └── router.go
│   └── usecase/               # Business logic
│       ├── auth_usecase.go
│       ├── category_usecase.go
│       └── event_usecase.go
├── migrations/                # Database migrations & seeders
│   ├── migrate.go
│   ├── seed_users.go
│   ├── seed_categories.go
│   └── seed_events.go
├── pkg/
│   ├── jwt/                  # JWT token handling
│   │   └── jwt.go
│   ├── response/             # API response helpers
│   │   └── response.go
│   ├── storage/              # File upload (S3)
│   │   └── s3_uploader.go
│   └── validator/            # Data validation
│       └── validator.go
├── public/                   # Static assets
│   └── uploads/
│       ├── events/
│       ├── users/
│       └── categories/
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── README.md                 # This file
```

## 🛠️ Useful Commands

### Database

```bash
# Run migrations
go run cmd/main.go -migrate

# Seed data
go run cmd/main.go -seed

# Fresh setup (drop & recreate)
go run cmd/main.go -fresh
```

### Development

```bash
# Run with hot reload (using air)
air

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code
golangci-lint run
```

### Deployment

```bash
# Build binary
make build

# Restart service
make restart

# View logs
make logs

# Full redeploy
make redeploy
```

## 🐳 Docker Setup (Optional)

If you want to use Docker:

```bash
# Build & run with docker-compose
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop containers
docker-compose down
```

## 🔐 Authentication

### JWT Token

Token is generated during login and must be included in the header for protected endpoints:

```bash
Authorization: Bearer <token>
```

### Seed Users

Default seeded users:

| Email | Password | Role |
|-------|----------|------|
| admin@gmail.com | admin123 | admin |
| promotor@gmail.com | promotor123 | promotor |
| user@gmail.com | user123 | user |

## 📤 Upload Assets

Assets are stored in the `public/uploads/` folder:

```
/uploads/events/
/uploads/users/
/uploads/categories/
```

Access via:
```
http://localhost:3000/uploads/events/filename.jpg
```

For production, use AWS S3 (set S3 config in `.env`).

## ⚙️ Troubleshooting

### Connection Refused
```
error connecting to database
```
**Solution:** Make sure PostgreSQL is running and config in `.env` is correct

### Port Already in Use
```
listen tcp :3000: bind: address already in use
```
**Solution:** Change `APP_PORT` in `.env` or kill the process using that port

### JWT Secret Not Set
```
JWT_SECRET is required
```
**Solution:** Set `JWT_SECRET` in `.env`

### S3 Connection Error
```
failed to initialize S3 uploader
```
**Solution:** Make sure S3 credentials in `.env` are correct or leave empty if not using S3

## 📖 Additional Documentation

- [API Documentation](https://2ikhh28mj3.apidog.io/)
- [Fiber Framework Docs](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/)
- [JWT Documentation](https://github.com/golang-jwt/jwt)

## 🤝 Contributing

Contributions are welcome! Please create a pull request with:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is a private project for Balai Budaya Citra (BCC).

## 📞 Support

For questions or assistance, please contact the development team.

---

**Last Updated:** March 14, 2026  
**Version:** 1.0.0
