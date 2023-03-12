package company

import "github.com/WeAreAmazingTeam/tcd-backend/user"

type (
	RequestCreateCompanyCashFlow struct {
		Status string `json:"status" binding:"required"`
		Amount int64  `json:"amount" binding:"required"`
		Note   string `json:"note" binding:"required"`
		User   user.User
	}

	RequestGetCompanyCashFlowByID struct {
		ID int `uri:"id" binding:"required"`
	}

	RequestDeleteCompanyCashFlow struct {
		User user.User
	}
)
