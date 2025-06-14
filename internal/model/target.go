package model

type Target struct {
	ID   	uint   `gorm:"primaryKey" json:"id"`
	Name 	string `json:"name"`
	Type 	string `json:"type"` // ping / http
	Address string `json:"address"`
}