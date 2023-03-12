package user

import (
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Register(req RequestRegister) (User, error)
	Login(req RequestLogin) (User, error)
	GetUserByID(ID int) (User, error)
	GetUserByEmail(Email string) (User, error)
	CheckDuplicateEmail(email string) (duplicate bool, err error)

	GetAllUser() ([]User, error)
	CreateUser(RequestCreateUser) (User, error)
	UpdateUser(RequestGetUserByID, RequestUpdateUser) (User, error)
	DeleteUser(RequestGetUserByID, RequestDeleteUser) (bool, error)

	GetWithdrawalRequestByID(id int) (UserWithdrawalRequest, error)
	CreateWithdrawalRequest(RequestCreateWithdrawalRequest) (UserWithdrawalRequest, error)
	UpdateUserWithdrawalRequest(RequestGetUserWithdrawalRequestByID, RequestUpdateUserWithdrawalRequest) (UserWithdrawalRequest, error)
	DeleteUserWithdrawalRequest(RequestGetUserWithdrawalRequestByID, RequestDeleteUserWithdrawalRequest) (bool, error)

	GetDataForgotPasswordByToken(token string) (UserForgotPasswordToken, error)
	CreateUserForgotPasswordToken(RequestCreateForgotPasswordToken) (UserForgotPasswordToken, error)
	DeleteForgotPasswordToken(UserForgotPasswordToken) (bool, error)

	AdminDataTablesUsers(*gin.Context) (helper.DataTables, error)

	AdminDataTablesWithdrawalRequest(*gin.Context) (helper.DataTables, error)

	UserDataTablesEMoneyFlow(*gin.Context, User) (helper.DataTables, error)
	UserDataTablesWithdrawalRequest(*gin.Context, User) (helper.DataTables, error)

	GetUserRegistered(condition string) (int, error)
	GetTotalWithdrawalRequest(condition string) (int, error)
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
