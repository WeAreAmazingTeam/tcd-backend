package company

import (
	"strconv"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
)

func (svc *service) CreateCompanyCashFlow(req RequestCreateCompanyCashFlow) (CompanyCashFlow, error) {
	companyCashFlow := CompanyCashFlow{}
	companyCashFlow.Status = req.Status
	companyCashFlow.Amount = req.Amount
	companyCashFlow.Note = req.Note
	companyCashFlow.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	companyCashFlowData, err := svc.repo.CreateCompanyCashFlow(companyCashFlow)

	if err != nil {
		return companyCashFlowData, err
	}

	return companyCashFlowData, nil
}

func (svc *service) DeleteCompanyCashFlow(reqDetail RequestGetCompanyCashFlowByID, reqDelete RequestDeleteCompanyCashFlow) (bool, error) {
	if constant.DELETED_BY {
		companyCashFlow, err := svc.repo.GetCompanyCashFlowByID(reqDetail.ID)

		if err != nil {
			return false, err
		}

		companyCashFlow.UpdatedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))
		companyCashFlow.DeletedAt = *helper.SetNowNT()
		companyCashFlow.DeletedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))

		status, err := svc.repo.DeleteCompanyCashFlow(companyCashFlow)

		if err != nil {
			return status, err
		}

		return status, nil
	}

	companyCashFlow := CompanyCashFlow{}
	companyCashFlow.ID = reqDetail.ID
	status, err := svc.repo.DeleteCompanyCashFlow(companyCashFlow)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) AdminDataTablesCompanyCashFlow(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesCompanyCashFlow, err := svc.repo.AdminDataTablesCompanyCashFlow(ctx)

	if err != nil {
		return dataTablesCompanyCashFlow, err
	}

	return dataTablesCompanyCashFlow, nil
}
