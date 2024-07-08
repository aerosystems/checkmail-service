package models

import (
	"time"
)

type Filter struct {
	Id           int    `gorm:"primaryKey;unique;autoIncrement"`
	Name         string ` gorm:"uniqueIndex:idx_name_coverage_project_token"`
	Type         string
	Coverage     string `gorm:"uniqueIndex:idx_name_coverage_project_token"`
	ProjectToken string `gorm:"uniqueIndex:idx_name_coverage_project_token"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
