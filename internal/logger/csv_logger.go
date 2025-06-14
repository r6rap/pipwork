package logger

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogEntry struct {
	Timestamp 	string
	Name 		string
	Type 		string
	Status 		string
	Latency 	string
	AverageLatency string
	Error 		string
}

func SaveLog(entry LogEntry) error {
	// get the current date in YYYY-MM-DD format
	dateStr := time.Now().Format("2006-01-02")

	// create the log file path: logs/YYY-MM-DD.csv
	fileName := filepath.Join("logs", dateStr+".csv")

	// check if the file already exists
	_, err := os.Stat(fileName)
	// determine if this is a new file
	newFile := os.IsNotExist(err)

	// open file for appending, create it if it doesn't exist, write-only mode
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file log: %v", err)
	}
	defer file.Close()

	// create a csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// if the file is new, write header
	if newFile {
		header := []string{"Timestamp", "Name", "Type", "Status", "Latency", "Error"}
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("failed to write header: %v", err)
		}
	}

	// prepare data log
	row := []string{
		entry.Timestamp,
		entry.Name,
		entry.Type,
		entry.Status,
		entry.Latency,
		entry.Error,
	}

	// write data log
	if err := writer.Write(row); err != nil {
		return fmt.Errorf("failed to write data log: %v", err)
	}

	return nil
}