package postgres

import (
	"fmt"
	// "hexagonalapp/app/domain/entity"
	// repo "hexagonalapp/app/domain/repository"
	// "hexagonalapp/app/services"

	"hexagonalapp/internal/platform/config"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
)

var DB *gorm.DB

// func Open(dsn string) (*gorm.DB, error) {
// 	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
// }

func DbConnect(cfg config.Config) *gorm.DB {
	var DBURL string
	if cfg.DBDriver == "mysql" {
		DBURL = cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"
	} else if cfg.DBDriver == "postgres" {
		DBURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName) //Build connection string
	} else if cfg.DBDriver == "supabase" {
		DBURL = cfg.Supabase //Build connection string
	}

	var logCtl logger.Interface

	switch cfg.DBDebugMode {
	case "DEBUG":
		if cfg.GormAdvancedLogger == "ENABLE" {
			logCtl = logger.Default.LogMode(logger.Info) // Debug için Info
		} else {
			logCtl = logger.Default.LogMode(logger.Error)
		}
	case "DEVELOPMENT":
		logCtl = logger.Default.LogMode(logger.Warn) // Test için Warn
	case "PRODUCTION":
		logCtl = logger.Default.LogMode(logger.Silent)
	default:
		logCtl = logger.Default.LogMode(logger.Error) // Default Error
	}

	db, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{
		Logger: logCtl,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // User → "user", Company → "company"
			TablePrefix:   "",
		},
	})
	db.Set("gorm:table_options", "charset=utf8")

	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	DB = db
	return db
}
