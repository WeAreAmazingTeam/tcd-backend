package handler

import (
	"fmt"
	"net/http"

	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/logs"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
)

type logsHandler struct {
	logsSvc logs.Service
}

func NewLogsHandler(logsService logs.Service) *logsHandler {
	return &logsHandler{logsSvc: logsService}
}

func (handler *logsHandler) AdminDataTablesActivityLogs(ctx *gin.Context) {
	dataTablesActivityLogs, err := handler.logsSvc.AdminDataTablesActivityLogs(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables activity logs failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesActivityLogs)
}

func (handler *logsHandler) AddLogsActivityAuth(ctx *gin.Context) {
	var req logs.RequestCreateActivityLog

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Add activity log failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userData := ctx.MustGet("userData").(user.User)
	response := helper.BasicAPIResponse(http.StatusCreated, "Add activity log successfully!")
	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v (from API request, by %v).", req.Content, userData.Name))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *logsHandler) AddLogsActivity(ctx *gin.Context) {
	var req logs.RequestCreateActivityLog

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Add activity log failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusCreated, "Add activity log successfully!")
	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v (from API request).", req.Content))

	ctx.JSON(http.StatusCreated, response)
}
