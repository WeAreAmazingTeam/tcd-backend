package logs

import (
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
)

func (svc *service) CreateActivityLog(ctx *gin.Context, content string) {
	log := ActivityLog{}
	log.Content = content
	log.UserAgent = ctx.Request.UserAgent()
	log.IpAddress = ctx.ClientIP()
	_, _ = svc.repo.SaveActivityLog(log)
}

func (svc *service) DeleteActivityLog(reqDetail RequestGetActivityLogByID, reqDelete RequestDeleteActivityLog) (bool, error) {
	log := ActivityLog{}
	log.ID = reqDetail.ID
	status, err := svc.repo.DeleteActivityLog(log)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) AdminDataTablesActivityLogs(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesActivityLogs, err := svc.repo.AdminDataTablesActivityLogs(ctx)

	if err != nil {
		return dataTablesActivityLogs, err
	}

	return dataTablesActivityLogs, nil
}

func (svc *service) CreateActivityWebhook(req RequestCreateActivityWebhook) {
	wh := ActivityWebhook{}
	wh.Endpoint = req.Endpoint
	wh.TriggeredFrom = req.TriggeredFrom
	wh.Properties = req.Properties
	_, _ = svc.repo.SaveActivityWebhook(wh)
}
