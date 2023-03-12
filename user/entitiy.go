package user

import (
	"time"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"
)

type (
	User struct {
		ID       int
		Role     string
		Name     string
		Email    string
		Password string
		EMoney   float64
		constant.CreatedUpdatedDeleted
	}

	UserEMoneyFlow struct {
		ID     int
		UserID int
		Status string
		Amount int64
		Note   string
		constant.CreatedUpdatedDeleted
	}

	UserWithdrawalRequest struct {
		ID     int
		UserID int
		Status string
		Amount int64
		Note   string
		constant.CreatedUpdatedDeleted
	}

	UserForgotPasswordToken struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func (UserEMoneyFlow) TableName() string {
	return "user_emoney_flow"
}
