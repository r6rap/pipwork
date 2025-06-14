package config

import (
	"github.com/r6rap/pipwork/internal/model"
	"github.com/r6rap/pipwork/internal/db"
)

func LoadTargets() ([]model.Target, error) {
	var targets []model.Target
	// find targets
	err := db.DB.Find(&targets).Error
	if err != nil {
		return nil, err
	}

	return targets, nil
}