package user

import (
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository interface {
	SaveUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByID(id int) (User, error)

	GetAllUser() ([]User, error)
	UpdateUser(User) (User, error)
	DeleteUser(User) (bool, error)

	GiveEMoneyToUser(userID, eMoney int) error

	CreateEMoneyFlow(UserEMoneyFlow) (UserEMoneyFlow, error)

	GetDataForgotPasswordByToken(token string) (UserForgotPasswordToken, error)
	CreateForgotPasswordToken(UserForgotPasswordToken) (UserForgotPasswordToken, error)
	DeleteForgotPasswordToken(UserForgotPasswordToken) (bool, error)

	GetWithdrawalRequestByID(id int) (UserWithdrawalRequest, error)
	CreateWithdrawalRequest(UserWithdrawalRequest) (UserWithdrawalRequest, error)
	UpdateUserWithdrawalRequest(UserWithdrawalRequest) (UserWithdrawalRequest, error)
	DeleteUserWithdrawalRequest(UserWithdrawalRequest) (bool, error)

	AdminDataTablesUsers(ctx *gin.Context) (helper.DataTables, error)
	AdminDataTablesWithdrawalRequest(*gin.Context) (helper.DataTables, error)

	UserDataTablesEMoneyFlow(*gin.Context, User) (helper.DataTables, error)
	UserDataTablesWithdrawalRequest(*gin.Context, User) (helper.DataTables, error)

	GetUserRegistered(condition string) (int, error)
	GetTotalWithdrawalRequest(condition string) (int, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{DB: db}
}
