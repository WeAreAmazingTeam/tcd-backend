package user

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/thanhpk/randstr"

	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (svc *service) Register(req RequestRegister) (User, error) {
	user := User{
		Name:  req.Name,
		Email: req.Email,
	}

	user.CreatedBy = helper.SetNS("New User")

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return user, err
	}

	user.Password = string(password)
	user.Role = "user"

	newUserData, err := svc.repo.SaveUser(user)

	if err != nil {
		return newUserData, err
	}

	return newUserData, nil
}

func (svc *service) Login(req RequestLogin) (User, error) {
	email := req.Email
	password := req.Password

	user, err := svc.repo.GetUserByEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("email not registered")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return user, err
	}

	return user, nil
}

func (svc *service) GetUserByID(id int) (User, error) {
	user, err := svc.repo.GetUserByID(id)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, fmt.Errorf("User with ID %d not found", id)
	}

	return user, nil
}

func (svc *service) CheckDuplicateEmail(email string) (duplicate bool, err error) {
	user, err := svc.repo.GetUserByEmail(email)

	if err != nil {
		if !helper.IsErrNoRows(err.Error()) {
			return false, err
		}
	}

	if user.ID != 0 {
		return true, nil
	}

	return false, nil
}

func (svc *service) GetAllUser() ([]User, error) {
	users, err := svc.repo.GetAllUser()

	if err != nil {
		return users, err
	}

	return users, nil
}

func (svc *service) CreateUser(req RequestCreateUser) (User, error) {
	user := User{}
	user.Role = req.Role
	user.Name = req.Name
	user.Email = req.Email
	user.EMoney = req.EMoney
	user.UpdatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return user, err
	}

	user.Password = string(password)

	newUserData, err := svc.repo.SaveUser(user)

	if err != nil {
		return newUserData, err
	}

	return newUserData, nil
}

func (svc *service) UpdateUserWithdrawalRequest(reqDetail RequestGetUserWithdrawalRequestByID, reqUpdate RequestUpdateUserWithdrawalRequest) (userWithdrawalRequest UserWithdrawalRequest, err error) {
	userWithdrawalRequest, err = svc.repo.GetWithdrawalRequestByID(reqDetail.ID)

	if err != nil {
		return userWithdrawalRequest, err
	}

	userWithdrawalRequest.Status = reqUpdate.Status
	userWithdrawalRequest.UpdatedBy = helper.SetNS(strconv.Itoa(reqUpdate.User.ID))

	updatedUserWithdrawalRequest, err := svc.repo.UpdateUserWithdrawalRequest(userWithdrawalRequest)

	if err != nil {
		return updatedUserWithdrawalRequest, err
	}

	if updatedUserWithdrawalRequest.Status == "approved" {
		svc.repo.CreateEMoneyFlow(UserEMoneyFlow{
			UserID: updatedUserWithdrawalRequest.UserID,
			Status: "out",
			Amount: int64(updatedUserWithdrawalRequest.Amount),
			Note:   "Withdrawal.",
		})

		withdrawalRequestUserData, err := svc.repo.GetUserByID(updatedUserWithdrawalRequest.UserID)

		if err != nil {
			return updatedUserWithdrawalRequest, err
		}

		withdrawalRequestUserData.EMoney = float64(withdrawalRequestUserData.EMoney) - float64(updatedUserWithdrawalRequest.Amount)

		if _, err := svc.repo.UpdateUser(withdrawalRequestUserData); err != nil {
			return updatedUserWithdrawalRequest, err
		}
	}

	return updatedUserWithdrawalRequest, nil
}

func (svc *service) DeleteUser(reqDetail RequestGetUserByID, reqDelete RequestDeleteUser) (bool, error) {
	if constant.DELETED_BY {
		user, err := svc.repo.GetUserByID(reqDetail.ID)

		if err != nil {
			return false, err
		}

		user.UpdatedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))
		user.DeletedAt = *helper.SetNowNT()
		user.DeletedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))

		status, err := svc.repo.DeleteUser(user)

		if err != nil {
			return status, err
		}

		return status, nil
	}

	user := User{}
	user.ID = reqDetail.ID
	status, err := svc.repo.DeleteUser(user)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) AdminDataTablesUsers(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesUsers, err := svc.repo.AdminDataTablesUsers(ctx)

	if err != nil {
		return dataTablesUsers, err
	}

	return dataTablesUsers, nil
}

func (svc *service) GetUserRegistered(condition string) (res int, err error) {
	res, err = svc.repo.GetUserRegistered(condition)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (svc *service) GetTotalWithdrawalRequest(condition string) (res int, err error) {
	res, err = svc.repo.GetTotalWithdrawalRequest(condition)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (svc *service) UserDataTablesEMoneyFlow(ctx *gin.Context, user User) (helper.DataTables, error) {
	dataTablesEMoneyFlow, err := svc.repo.UserDataTablesEMoneyFlow(ctx, user)

	if err != nil {
		return dataTablesEMoneyFlow, err
	}

	return dataTablesEMoneyFlow, nil
}

func (svc *service) CreateWithdrawalRequest(req RequestCreateWithdrawalRequest) (UserWithdrawalRequest, error) {
	withdrawalRequest := UserWithdrawalRequest{}
	withdrawalRequest.UserID = req.User.ID
	withdrawalRequest.Status = "pending"
	withdrawalRequest.Amount = req.Amount
	withdrawalRequest.Note = req.Note
	withdrawalRequest.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	withdrawalRequestData, err := svc.repo.CreateWithdrawalRequest(withdrawalRequest)

	if err != nil {
		return withdrawalRequestData, err
	}

	return withdrawalRequestData, nil
}

func (svc *service) UserDataTablesWithdrawalRequest(ctx *gin.Context, user User) (helper.DataTables, error) {
	dataTablesWithdrawalRequest, err := svc.repo.UserDataTablesWithdrawalRequest(ctx, user)

	if err != nil {
		return dataTablesWithdrawalRequest, err
	}

	return dataTablesWithdrawalRequest, nil
}

func (svc *service) AdminDataTablesWithdrawalRequest(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesWithdrawalRequest, err := svc.repo.AdminDataTablesWithdrawalRequest(ctx)

	if err != nil {
		return dataTablesWithdrawalRequest, err
	}

	return dataTablesWithdrawalRequest, nil
}

func (svc *service) GetWithdrawalRequestByID(id int) (UserWithdrawalRequest, error) {
	userWithdrawalRequest, err := svc.repo.GetWithdrawalRequestByID(id)

	if err != nil {
		return userWithdrawalRequest, err
	}

	if userWithdrawalRequest.ID == 0 {
		return userWithdrawalRequest, fmt.Errorf("withdrawal request with ID %d not found", id)
	}

	return userWithdrawalRequest, nil
}

func (svc *service) UpdateUser(reqDetail RequestGetUserByID, reqUpdate RequestUpdateUser) (user User, err error) {
	user, err = svc.repo.GetUserByID(reqDetail.ID)

	if err != nil {
		return user, err
	}

	user.ID = reqDetail.ID
	user.Role = reqUpdate.Role
	user.Name = reqUpdate.Name
	user.Email = reqUpdate.Email
	user.EMoney = reqUpdate.EMoney
	user.UpdatedBy = helper.SetNS(strconv.Itoa(reqUpdate.User.ID))

	if reqUpdate.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(reqUpdate.Password), bcrypt.DefaultCost)

		if err != nil {
			return user, err
		}

		user.Password = string(password)
	} else {
		existingUser, err := svc.GetUserByID(user.ID)

		if err != nil {
			return user, err
		}

		user.Password = existingUser.Password
	}

	updatedUser, err := svc.repo.UpdateUser(user)

	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (svc *service) DeleteUserWithdrawalRequest(reqDetail RequestGetUserWithdrawalRequestByID, reqDelete RequestDeleteUserWithdrawalRequest) (bool, error) {
	if constant.DELETED_BY {
		userWithdrawalRequest, err := svc.repo.GetWithdrawalRequestByID(reqDetail.ID)

		if err != nil {
			return false, err
		}

		userWithdrawalRequest.UpdatedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))
		userWithdrawalRequest.DeletedAt = *helper.SetNowNT()
		userWithdrawalRequest.DeletedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))

		status, err := svc.repo.DeleteUserWithdrawalRequest(userWithdrawalRequest)

		if err != nil {
			return status, err
		}

		return status, nil
	}

	userWithdrawalRequest := UserWithdrawalRequest{}
	userWithdrawalRequest.ID = reqDetail.ID
	status, err := svc.repo.DeleteUserWithdrawalRequest(userWithdrawalRequest)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) CreateUserForgotPasswordToken(req RequestCreateForgotPasswordToken) (UserForgotPasswordToken, error) {
	userForgotPasswordToken := UserForgotPasswordToken{}
	userForgotPasswordToken.UserID = req.User.ID
	userForgotPasswordToken.Token = randstr.String(69)

	userForgotPasswordTokenData, err := svc.repo.CreateForgotPasswordToken(userForgotPasswordToken)

	if err != nil {
		return userForgotPasswordTokenData, err
	}

	return userForgotPasswordTokenData, nil
}

func (svc *service) GetUserByEmail(email string) (User, error) {
	user, err := svc.repo.GetUserByEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, fmt.Errorf("User with email %v not found", email)
	}

	return user, nil
}

func (svc *service) GetDataForgotPasswordByToken(token string) (UserForgotPasswordToken, error) {
	userForgotPasswordToken, err := svc.repo.GetDataForgotPasswordByToken(token)

	if err != nil {
		return userForgotPasswordToken, err
	}

	if userForgotPasswordToken.ID == 0 {
		return userForgotPasswordToken, fmt.Errorf("get data forgot password with token %s not found", token)
	}

	return userForgotPasswordToken, nil
}

func (svc *service) DeleteForgotPasswordToken(userForgotPasswordToken UserForgotPasswordToken) (bool, error) {
	status, err := svc.repo.DeleteForgotPasswordToken(userForgotPasswordToken)

	if err != nil {
		return status, err
	}

	return status, nil
}
