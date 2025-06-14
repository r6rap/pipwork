package model

import "time"

type MonitoringLog struct {
	ID 				uint		`gorm:"primaryKey"`
	Timestamp 		time.Time
	Name 			string
	Type 			string
	Status 			string
	Latency 		string
	AverageLatency 	string
	Error 			string
}