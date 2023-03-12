package company

import (
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateCompanyCashFlow(RequestCreateCompanyCashFlow) (CompanyCashFlow, error)
	DeleteCompanyCashFlow(RequestGetCompanyCashFlowByID, RequestDeleteCompanyCashFlow) (bool, error)

	AdminDataTablesCompanyCashFlow(*gin.Context) (helper.DataTables, error)
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
