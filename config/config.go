package config

import "gorm.io/gorm"

// Config stores application configurations.
type Config struct {
	DB *gorm.DB
}

// NewConfig creates Config.
func NewConfig() (*Config, error) {
	db, err := ConnDB()
	if err != nil {
		return nil, err
	}
	return &Config{DB: db}, nil
}
