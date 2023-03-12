package transaction

import (
	"github.com/WeAreAmazingTeam/tcd-backend/campaign"
	"github.com/WeAreAmazingTeam/tcd-backend/company"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/payment"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
)

type Service interface {
	GetAllTransaction(*gin.Context) ([]Transaction, error)
	GetTransactionByID(RequestGetTransactionByID) (Transaction, error)
	GetTransactionByCampaignID(*gin.Context, RequestGetTransactionByCampaignID) ([]TransactionWithUserName, error)
	GetTransactionByUserID(*gin.Context, RequestGetTransactionByUserID) ([]Transaction, error)
	CreateTransaction(req RequestCreateTransaction, campaignName string) (Transaction, error)
	CreateTransactionWithEMoney(req RequestCreateTransactionWithEMoney, campaignName string) (Transaction, error)
	CreateAnonymousTransaction(req RequestCreateAnonymousTransaction, campaignName string) (Transaction, error)
	DeleteTransaction(RequestGetTransactionByID, RequestDeleteTransaction) (bool, error)

	AdminDataTablesTransactions(*gin.Context) (helper.DataTables, error)
	UserDataTablesTransactions(*gin.Context, user.User) (helper.DataTables, error)

	ProcessRequestFromMidtrans(req MidtransRequest) error

	GetTotalTransaction(condition string) (int, error)
}

type service struct {
	repo         Repository
	campaignRepo campaign.Repository
	userRepo     user.Repository
	companyRepo  company.Repository
	campaignSvc  campaign.Service
	paymentSvc   payment.Service
}

func NewService(
	repository Repository,
	campaignRepository campaign.Repository,
	userRepository user.Repository,
	companyRepository company.Repository,
	campaignService campaign.Service,
	paymentService payment.Service,
) *service {
	return &service{
		repo:         repository,
		campaignRepo: campaignRepository,
		userRepo:     userRepository,
		companyRepo:  companyRepository,
		campaignSvc:  campaignService,
		paymentSvc:   paymentService,
	}
}
