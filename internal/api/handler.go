package api

import (
	"encoding/json"
	"net/http"
	// "strconv"

	"pipwork/internal/db"
	"pipwork/internal/model"

	"github.com/gorilla/mux"
)

func GetLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logs []model.MonitoringLog

	page := 1
	limit := 10

	query := db.DB

	offset := (page - 1) * limit

	// set response header
	w.Header().Set("Content-Type", "application/json")

	// filter by "name" query parameter if present
	name := r.URL.Query().Get("name")
	if name != "" {
		query = query.Where("name = ?", name)
	}

	var total int64

	// retrieve logs ordered by timestamp in descending order
	err := query.Order("timestamp desc").Limit(limit).Offset(offset).Find(&logs).Count(&total).Error
	if err != nil {
		http.Error(w, "Failed to get data logs", http.StatusInternalServerError)
		return
	}

	if total > int64(limit) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": logs,
			"page": page,
			"limit": limit,
			"total": total,
			"total pages": int((total + int64(limit) - 1) / int64(limit)),
		})
	} else {
		json.NewEncoder(w).Encode(logs)
	}
}

func GetTarget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var targets []model.Target

	err := db.DB.Find(&targets).Error
	if err != nil {
		http.Error(w, "Failed to get targets", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(targets)
}

func GetTargetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	var target model.Target

	// get first row
	err := db.DB.First(&target, params["id"]).Error
	if err != nil {
		http.Error(w, "Target not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(target)
}

func CreateTarget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var target model.Target

	// decode JSON request body
	err := json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// validate required fields
	if target.Name == "" || target.Type == "" || target.Address == "" {
		http.Error(w, "Fields cannot be empty", http.StatusBadRequest)
		return
	}

	// insert target into the database
	err = db.DB.Create(&target).Error
	if err != nil {
		http.Error(w, "Failed to save target", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully added new target",
	})
}

func UpdateTarget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	var target model.Target

	err := db.DB.First(&target, params["id"]).Error
	if err != nil {
		http.Error(w, "Target not found", http.StatusNotFound)
		return
	}

	var updated model.Target
	err = json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusInternalServerError)
		return
	}

	target.Name = updated.Name
	target.Type = updated.Type
	target.Address = updated.Address

	db.DB.Save(&target)
	json.NewEncoder(w).Encode(target)
}

func DeleteTarget(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := mux.Vars(r)
	err := db.DB.Delete(&model.Target{}, params["id"]).Error
	if err != nil {
		http.Error(w, "Failed to delete target", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
