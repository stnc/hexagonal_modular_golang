package config

import "os"

type Config struct {
	App                string
	AppPort            string
	DBName             string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBDriver           string
	MongoDBURI         string
	MongoDBName        string
	MongoDBUser        string
	MongoDBPassword    string
	DBSSLMode          string
	RedisHost          string
	RedisPort          string
	RedisPassword      string
	RedisDB            int
	DBDebugMode        string
	GormAdvancedLogger string
	Supabase           string
	EnvName            string
	// PostgresDSN        string

}

func Load() Config {
	return Config{
		App:         getEnv("APP", "all"),
		AppPort:     getEnv("APP_PORT", "8080"),
		DBName:      getEnv("DB_NAME", "langCamp"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "changeme"),
		DBDriver:    getEnv("DBDriver", "postgres"),
		MongoDBURI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGODB_DB", "admin"),
		MongoDBUser: getEnv("MONGODB_DB_USER", "root"),
		MongoDBPassword:    getEnv("MONGODB_DB_PASSWORD", "examplepassword"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		RedisHost:          getEnv("REDIS_HOST", "redis"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            0,
		DBDebugMode:        getEnv("DB_DEBUG_MODE", "DEBUG"),
		GormAdvancedLogger: getEnv("GORM_ZAP_LOGGER", "ENABLE"),
		Supabase:           getEnv("SUPABASE_DSN", "s"),
		EnvName:            getEnv("APP_ENV", "development"),
		// PostgresDSN:        getEnv("POSTGRES_DSN", "host=localhost user=postgres password=postgres dbname=hexagonalapp port=5432 sslmode=disable TimeZone=UTC"),

	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
