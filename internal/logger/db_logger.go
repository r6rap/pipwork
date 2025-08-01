package logger

import (
	"time"

	"github.com/r6rap/pipwork/internal/db"
	"github.com/r6rap/pipwork/internal/model"
)

func SaveToDB(entry LogEntry) error {
	log := model.MonitoringLog {
		Timestamp: time.Now(),
		Name: entry.Name,
		Type: entry.Type,
		Status: entry.Status,
		Latency: entry.Latency,
		AverageLatency: entry.AverageLatency,
		Error: entry.Error,
	}

	// insert data to table
	return db.DB.Create(&log).Error
}