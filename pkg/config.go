package pkg

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

// DatabaseConfig holds MySQL database configuration
type DatabaseConfig struct {
	Host               string
	Port               string
	Username           string
	Password           string
	Database           string
	MaxIdleConnections int
	MaxOpenConnections int
	MaxLifetime        time.Duration
	SSLMode            string
}

// ServerConfig holds all server-related configuration
type ServerConfig struct {
	Host             string
	Port             int
	Environment      string
	LogLevel         string
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	EnvPath          string
}

// Config is the root configuration structure
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

func (config *Config) Validate() {
	if err := config.Database.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := config.Server.Validate(); err != nil {
		log.Fatal(err)
	}
}

func (serverConfig *ServerConfig) Validate() error {
	return nil
}

func (serverConfig *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
}

// DefaultServerConfig returns a ServerConfig with default values
func getServerConfig() ServerConfig {
	return ServerConfig{
		Host:        GetEnvOrDefault("SERVER_HOST", "0.0.0.0"),
		Port:        GetEnvAsIntOrDefault("SERVER_PORT", 8080),
		Environment: GetEnvOrDefault("ENVIRONMENT", "development"),
		LogLevel:    GetEnvOrDefault("LOG_LEVEL", "info"),
		AllowedOrigins: GetEnvAsSliceOrDefault(
			"ALLOWED_ORIGINS",
			[]string{"http://localhost:3000"},
		),
		AllowedMethods: GetEnvAsSliceOrDefault(
			"ALLOWED_METHODS",
			[]string{"GET", "POST", "PUT", "DELETE"},
		),
		AllowedHeaders: GetEnvAsSliceOrDefault(
			"ALLOWED_HEADERS",
			defaultAllowedHeaders(),
		),
		AllowCredentials: GetEnvAsBoolOrDefault(
			"ALLOW_CREDENTIALS",
			false,
		),
		ExposedHeaders: GetEnvAsSliceOrDefault(
			"EXPOSED_HEADERS",
			defaultExposedHeaders(),
		),
		EnvPath: ".env",
	}
}

// New creates a new Config instance
func New(envPath string) (*Config, error) {
	if err := loadEnv(envPath); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	return &Config{
		Server:   getServerConfig(),
		Database: getDatabaseConfig(),
	}, nil
}

// DefaultDatabaseConfig returns default database configuration from environment variables
func getDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:               GetEnvOrDefault("MYSQL_DATABASE_HOST", "localhost"),
		Port:               GetEnvOrDefault("MYSQL_DATABASE_PORT", "3306"),
		Username:           GetEnvOrDefault("MYSQL_USERNAME", "root"),
		Password:           GetEnvOrDefault("MYSQL_PASSWORD", ""),
		Database:           GetEnvOrDefault("MYSQL_DATABASE_NAME", "haven"),
		SSLMode:            GetEnvOrDefault("MYSQL_SSL_MODE", "disable"),
		MaxIdleConnections: GetEnvAsIntOrDefault("MYSQL_MAX_IDLE_CONNS", 10),
		MaxOpenConnections: GetEnvAsIntOrDefault("MYSQL_MAX_OPEN_CONNS", 100),
		MaxLifetime:        time.Duration(GetEnvAsIntOrDefault("MYSQL_CONN_MAX_LIFETIME", 3600)) * time.Second,
	}
}

// ConnectionString generates MySQL connection string
func (databaseConfig DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		databaseConfig.Username,
		databaseConfig.Password,
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Database,
	)
}

// GetDSN returns the database connection string
func (databaseConfig DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		databaseConfig.Host, databaseConfig.Port, databaseConfig.Username, databaseConfig.Password, databaseConfig.Database, databaseConfig.SSLMode,
	)
}

// GetDB initializes and returns a configured GORM DB instance
func (databaseConfig DatabaseConfig) GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(databaseConfig.ConnectionString()), &gorm.Config{
		Logger: DatabaseQueryLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(databaseConfig.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(databaseConfig.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(databaseConfig.MaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Validate checks if the database configuration is valid
func (databaseConfig DatabaseConfig) Validate() error {
	if databaseConfig.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}

	if databaseConfig.Port == "" {
		return fmt.Errorf("database port cannot be empty")
	}

	if port, err := strconv.Atoi(databaseConfig.Port); err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid database port: %s", databaseConfig.Port)
	}

	if databaseConfig.Username == "" {
		return fmt.Errorf("database username cannot be empty")
	}

	if databaseConfig.Database == "" {
		return fmt.Errorf("database name cannot be empty")
	}

	return nil
}

func loadEnv(envPath string) error {
	if envPath == "" {
		envPath = ".env"
	}

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil // Skip if .env doesn't exist
	}

	return godotenv.Load(envPath)
}

func defaultAllowedHeaders() []string {
	return []string{
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"accept",
		"origin",
		"Cache-Control",
		"X-Requested-With",
		"Access-Control-Allow-Origin",
	}
}

func defaultExposedHeaders() []string {
	return []string{
		"Content-Length",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"Content-Type",
	}
}
