package logs

import (
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateActivityLog(*gin.Context, string)
	DeleteActivityLog(RequestGetActivityLogByID, RequestDeleteActivityLog) (bool, error)

	CreateActivityWebhook(RequestCreateActivityWebhook)

	AdminDataTablesActivityLogs(*gin.Context) (helper.DataTables, error)
}

type service struct {
	repo Repository
}

func NewService(
	repository Repository,
) *service {
	return &service{
		repo: repository,
	}
}
