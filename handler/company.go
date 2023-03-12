package handler

import (
	"fmt"
	"net/http"

	"github.com/WeAreAmazingTeam/tcd-backend/company"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/logs"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
)

type companyHandler struct {
	companySvc company.Service
	logsSvc    logs.Service
}

func NewCompanyHandler(
	companyService company.Service,
	logsService logs.Service,
) *companyHandler {
	return &companyHandler{
		companySvc: companyService,
		logsSvc:    logsService,
	}
}

func (handler *companyHandler) AdminDataTablesCompanyCashFlow(ctx *gin.Context) {
	dataTablesCompanyCashFlow, err := handler.companySvc.AdminDataTablesCompanyCashFlow(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables company cash flow request failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesCompanyCashFlow)
}

func (handler *companyHandler) CreateCompanyCashFlow(ctx *gin.Context) {
	var req company.RequestCreateCompanyCashFlow

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Create company cash flow failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	req.User = ctx.MustGet("userData").(user.User)

	companyCashFlowData, err := handler.companySvc.CreateCompanyCashFlow(req)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Create company cash flow failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := company.FormatCompanyCashFlowData(companyCashFlowData)
	response := helper.APIResponse(http.StatusCreated, "Create company cash flow successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v creating company cash flow id %v.", req.User.Name, companyCashFlowData.ID))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *companyHandler) DeleteCompanyCashFlow(ctx *gin.Context) {
	var reqID company.RequestGetCompanyCashFlowByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete company cash flow failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete company.RequestDeleteCompanyCashFlow

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete company cash flow failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.companySvc.DeleteCompanyCashFlow(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete company cash flow failed!", fmt.Sprintf("Company Cash Flow with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete company cash flow failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete company cash flow successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting company cash flow id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}
