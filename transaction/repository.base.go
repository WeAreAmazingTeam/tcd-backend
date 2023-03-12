package transaction

import (
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllTransaction(*gin.Context) ([]Transaction, error)
	GetTransactionByCampaignId(ctx *gin.Context, campaignID int) ([]TransactionWithUserName, error)
	GetTransactionByUserID(ctx *gin.Context, userID int) ([]Transaction, error)
	GetTransactionByID(id int) (Transaction, error)
	GetTransactionByCode(code string) (Transaction, error)
	SaveTransaction(transaction Transaction) (Transaction, error)
	UpdateTransaction(Transaction) (Transaction, error)
	DeleteTransaction(Transaction) (bool, error)

	AdminDataTablesTransactions(*gin.Context) (helper.DataTables, error)
	UserDataTablesTransactions(*gin.Context, user.User) (helper.DataTables, error)

	GetTotalTransaction(condition string) (int, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{DB: db}
}
