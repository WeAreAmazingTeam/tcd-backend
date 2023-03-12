package logs

import (
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository interface {
	SaveActivityLog(ActivityLog) (ActivityLog, error)
	DeleteActivityLog(ActivityLog) (bool, error)

	SaveActivityWebhook(ActivityWebhook) (ActivityWebhook, error)

	AdminDataTablesActivityLogs(ctx *gin.Context) (helper.DataTables, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{DB: db}
}
