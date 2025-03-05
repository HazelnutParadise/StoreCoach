package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	host     = "localhost"
	port     = "9920"
	user     = "user"
	password = "password"
	dbname   = "StoreCoach"
	sslmode  = "disable"
)

func init() {
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode + " TimeZone=Asia/Taipei"
	var err error
	for err == nil {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

}
