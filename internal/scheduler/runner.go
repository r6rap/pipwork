package scheduler

import (
	"fmt"
	"log"
	"github.com/r6rap/pipwork/internal/checker"
	"github.com/r6rap/pipwork/internal/db"
	"github.com/r6rap/pipwork/internal/logger"
	"github.com/r6rap/pipwork/internal/model"
	"time"
)

type FinalResult struct{
	Name string
	Status string
	Latency time.Duration
	Avg time.Duration
	StatusCode int
	SSLValidation bool
	ExpiredTime time.Time
	DaysLeft int
	Error string
}

func StartMonitoring(targets []model.Target, interval time.Duration) {
	// create a ticker that sends a signal at the specified interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Println("🚀 Monitoring dimulai setiap", interval)

	for{
		select{
		// wait for the next tick (interval reached)
		case <-ticker.C:
			fmt.Println("\n⏱️  waktu:", time.Now().Format("15:04:05"))
			// launch a goroutine for each target to run concurrently
			for _, t := range targets{
				go monitorTarget(t)
			}
		}
	}
}

func monitorTarget(t model.Target) {
	now := time.Now().Format("2006-01-02 15:04:05")

	switch t.Type {
		case "http":
			res := checker.CheckHTTP(t.Address)
			ssl := checker.CheckSSL(t.Address)
			avg := countAverage(t.Name)
			
			result := FinalResult{
				Name: t.Name,
				Status: res.Status,
				Latency: res.Latency,
				Avg: time.Duration(avg),
				StatusCode: res.StatusCode,
				SSLValidation: ssl.Valid,
				ExpiredTime: ssl.ExpiredTime,
				DaysLeft: ssl.DaysUntilExpiry,
				Error: combineErrors(res.Error, ssl.Error),
			}

			fmt.Println("[HTTP]", result.Name+":")
			fmt.Printf("Status: %s, Latency: %s (%d)\nPERFORMANCE: Avg: %v\nSSL Validation: %t, Expired Date: %s, Days left: %d\nError: %s\n",
				result.Status,
				result.Latency,
				result.StatusCode,
				result.Avg,
				result.SSLValidation,
				result.ExpiredTime.Format("2006-01-02"),
				result.DaysLeft,
				result.Error,
			)


			logger.SaveToDB(logger.LogEntry{
				Timestamp: now,
				Name:      t.Name,
				Type:      "http",
				Status:    res.Status,
				Latency:   res.Latency.String(),
				Error:     combineErrors(res.Error, ssl.Error),
			})
		case "ping":
			res := checker.CheckPING(t.Address)
			fmt.Printf("[PING] %s: %s (%s)\n", t.Name, res.Status, res.Latency)
			logger.SaveToDB(logger.LogEntry{
				Timestamp: now,
				Name:      t.Name,
				Type:      "ping",
				Status:    res.Status,
				Latency:   res.Latency.String(),
				Error:     res.Error,
			})
		// handle unknown target types
		default:
			fmt.Printf("[UNKNOWN] %s: Tipe tidak dikenal: %s\n", t.Name, t.Type)
	}
}

func combineErrors(err1, err2 string) string {
	if err1 != "" && err2 != "" {
		return err1 +" | "+ err2
	}

	return err1 + err2
}

func countAverage(name string) time.Duration {
	var avg []string

	// custom time to count average
	from := time.Now().Add(-5 * time.Hour)
	to := time.Now()

	query := db.DB

	err := query.Table("monitoring_logs").Select("latency").Where("name = ? AND timestamp BETWEEN ? AND ?", name, from, to).Order("timestamp desc").Pluck("latency", &avg).Error
	if err != nil {
		log.Println("Error counting average:", err)
		return 0
	}

	if len(avg) == 0 {
		return 0
	}

	var total time.Duration

	for _, s := range avg {
		val, err := time.ParseDuration(s)
		if err != nil {
			log.Println("Error parsing", err)
			continue
		}

		total += val
	}

	return total / time.Duration(len(avg))
}