package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"pipwork/internal/model"
)

var DB *gorm.DB

func ConnectMySQL() {
	dsn := "root:@tcp(127.0.0.1:3306)/network_monitor?charset=utf8mb4&parseTime=True&loc=Asia%2FJakarta"
	var err error

	// open MySQL connection using GORM
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// get the underlying *sql.DB to configure connection pool settings
	sqlDB, _ := DB.DB()
	// set maximum lifetime of a connection
	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	// set maximum number of open connections
	sqlDB.SetMaxOpenConns(10)
	// set maximum number of idle connections
	sqlDB.SetMaxIdleConns(5)

	// auto migrate schema for MonitoringLog model
	err = DB.AutoMigrate(&model.MonitoringLog{}, &model.Target{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	fmt.Println("connected to MySQL & migration successful")
}