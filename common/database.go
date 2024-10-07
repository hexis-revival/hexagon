package common

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfiguration struct {
	Host        string
	Port        int
	Username    string
	Password    string
	Database    string
	MaxIdle     int
	MaxOpen     int
	MaxLifetime time.Duration
}

func (config *DatabaseConfiguration) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d TimeZone=UTC",
		config.Host, config.Username, config.Password, config.Database, config.Port,
	)
}

func CreateDatabaseSession(config *DatabaseConfiguration) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(config.MaxIdle)
	sqlDB.SetMaxOpenConns(config.MaxOpen)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)

	return db, nil
}

func DoDatabaseMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Stats{})
}
