package config

import "gorm.io/gorm"

// Config stores application configurations.
type Config struct {
	// Database connection.
	DB *gorm.DB

	// Limitation of the number of articles in a page. Default is 6.
	LimitInPage int
}

// NewConfig creates Config.
func NewConfig() (*Config, error) {
	db, err := ConnDB()
	if err != nil {
		return nil, err
	}
	return &Config{DB: db, LimitInPage: 6}, nil
}
