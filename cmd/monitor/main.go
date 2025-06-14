package main

import (
	"time"
	"log"

	"pipwork/internal/config"
	"pipwork/internal/scheduler"
	"pipwork/internal/db"
)

func main() {
	db.ConnectMySQL()

	target, err := config.LoadTargets()
	if err != nil {
		log.Fatalf("failed load target: %v", err)
	}

	scheduler.StartMonitoring(target, 1*time.Minute)
}