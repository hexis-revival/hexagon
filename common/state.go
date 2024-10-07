package common

import (
	"gorm.io/gorm"
)

type StateConfiguration struct {
	Database *DatabaseConfiguration
	DataPath string
}

func NewStateConfiguration() *StateConfiguration {
	return &StateConfiguration{
		Database: &DatabaseConfiguration{},
		DataPath: ".data",
	}
}

type State struct {
	Database *gorm.DB
	Storage  Storage
}

func NewState(config *StateConfiguration) (*State, error) {
	db, err := CreateDatabaseSession(config.Database)
	if err != nil {
		return nil, err
	}

	storage := NewFileStorage(config.DataPath)

	return &State{
		Database: db,
		Storage:  storage,
	}, nil
}
