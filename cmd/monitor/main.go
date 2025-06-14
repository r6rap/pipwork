package main

import (
	"time"
	"log"

	"github.com/r6rap/pipwork/internal/config"
	"github.com/r6rap/pipwork/internal/scheduler"
	"github.com/r6rap/pipwork/internal/db"
)

func main() {
	db.ConnectMySQL()

	target, err := config.LoadTargets()
	if err != nil {
		log.Fatalf("failed load target: %v", err)
	}

	scheduler.StartMonitoring(target, 1*time.Minute)
}