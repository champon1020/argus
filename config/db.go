package config

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True",
		os.Getenv("ARGUS_DB_USER"),
		os.Getenv("ARGUS_DB_PASSWORD"),
		os.Getenv("ARGUS_DB_HOST"),
		os.Getenv("ARGUS_DB_PORT"),
		os.Getenv("ARGUS_DB_NAME"),
	)
}

// ConnDB conntects to database.
func ConnDB() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn()), &gorm.Config{})
}
