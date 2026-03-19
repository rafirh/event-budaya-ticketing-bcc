package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName        string
	AppEnv         string
	AppPort        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpiryHours int
	S3Key          string
	S3Secret       string
	S3Bucket       string
	S3Region       string
	S3PublicBase   string
	MidtransEnv    string
	MidtransServer string
	MidtransClient string
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	AppConfig = &Config{
		AppName:        getEnv("APP_NAME", "event-budaya-ticketing"),
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnv("APP_PORT", "3000"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "secret"),
		DBName:         getEnv("DB_NAME", "event_budaya_db"),
		JWTSecret:      getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		JWTExpiryHours: jwtExpiry,
		S3Key:          getEnv("S3_KEY", ""),
		S3Secret:       getEnv("S3_SECRET", ""),
		S3Bucket:       getEnv("S3_BUCKET", ""),
		S3Region:       getEnv("S3_REGION", ""),
		S3PublicBase:   getEnv("S3_PUBLIC_BASE_URL", ""),
		MidtransEnv:    getEnv("MIDTRANS_ENV", "sandbox"),
		MidtransServer: getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClient: getEnv("MIDTRANS_CLIENT_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
