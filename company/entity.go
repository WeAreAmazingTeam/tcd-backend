package company

import "github.com/WeAreAmazingTeam/tcd-backend/constant"

type (
	CompanyCashFlow struct {
		ID     int
		Status string
		Amount int64
		Note   string
		constant.CreatedUpdatedDeleted
	}
)

func (CompanyCashFlow) TableName() string {
	return "company_cash_flow"
}
