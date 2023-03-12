package chart

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetChart(chartName, year string) (Chart, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{DB: db}
}
