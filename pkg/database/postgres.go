package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// PostGIS 확장 설치 확인
	var result string
	if err := db.Raw("SELECT PostGIS_version()").Scan(&result).Error; err != nil {
		log.Println("Warning: PostGIS not installed or not available")
	} else {
		log.Printf("PostGIS version: %s", result)
	}

	// 연결 풀 설정
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connection established")
	return db, nil
}
