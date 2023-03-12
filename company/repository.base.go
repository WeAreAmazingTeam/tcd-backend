package company

import (
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository interface {
	GetCompanyCashFlowByID(id int) (CompanyCashFlow, error)
	CreateCompanyCashFlow(CompanyCashFlow) (CompanyCashFlow, error)
	DeleteCompanyCashFlow(CompanyCashFlow) (bool, error)

	AdminDataTablesCompanyCashFlow(*gin.Context) (helper.DataTables, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{DB: db}
}
