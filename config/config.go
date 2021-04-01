package config

import "gorm.io/gorm"

// Config stores application configurations.
type Config struct {
	// Database connection.
	DB *gorm.DB

	// Limitation on the number of articles in a page. Default is 6.
	LimitOnNumArticles int

	// Limitation on the number of private articles in a page. Default is 10.
	LimitOnNumPrivArticles int

	// Limitations on the number of images in a page. Default is 12.
	LimitOnNumImages int

	// Bucket name of Google Cloud Storage.
	StorageBucketName string
}

// NewConfig creates Config.
func NewConfig() (*Config, error) {
	db, err := ConnDB()
	if err != nil {
		return nil, err
	}

	config := &Config{
		DB:                     db,
		LimitOnNumArticles:     6,
		LimitOnNumPrivArticles: 10,
		LimitOnNumImages:       12,
		StorageBucketName:      "myblog-argus",
	}

	return config, nil
}
