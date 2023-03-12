package logs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
)

func (repo *repository) SaveActivityLog(activityLog ActivityLog) (ActivityLog, error) {
	if err := repo.DB.Create(&activityLog).Error; err != nil {
		return activityLog, err
	}
	return activityLog, nil
}

func (repo *repository) DeleteActivityLog(activityLog ActivityLog) (bool, error) {
	if err := repo.DB.Delete(&activityLog).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) AdminDataTablesActivityLogs(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesActivityLogs
		likes string = ""
		order string = ""
		limit string = ""
	)

	var (
		no       int = 1
		total    int = 0
		filtered int = 0
	)

	var data []map[string]any

	listOrder := []string{"", "content", "user_agent", "ip_address", "created_at", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf(`(
			content LIKE '%%%s%%' OR user_agent LIKE '%%%s%%' OR ip_address LIKE '%%%s%%'
		)`, searchValue, searchValue, searchValue)
	}

	if orderColumn != "" {
		orderType := ctx.Query("order[0][dir]")
		orderColumn, _ := strconv.Atoi(orderColumn)
		order = fmt.Sprintf("ORDER BY %s %s", listOrder[orderColumn], strings.ToUpper(orderType))
	} else {
		order = "ORDER BY id DESC"
	}

	if starting != -1 {
		length, _ := strconv.Atoi(ctx.Query("length"))
		limit = fmt.Sprintf("LIMIT %v OFFSET %v", length, starting)
		no = starting + 1
	}

	if likes != "" {
		query = fmt.Sprintf("%s AND %s", query, likes)

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesActivityLogs), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesActivityLogs)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesActivityLogs)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		total = filtered
	}

	if order != "" {
		query = fmt.Sprintf("%s %s", query, order)
	}

	query = fmt.Sprintf("%s %s", query, limit)

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(query)).Rows()

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := ActivityLog{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Content,
			&tmp.UserAgent,
			&tmp.IpAddress,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return result, err
		}

		data = append(data, map[string]any{
			"no":         no,
			"id":         tmp.ID,
			"content":    tmp.Content,
			"user_agent": tmp.UserAgent,
			"ip_address": tmp.IpAddress,
			"created_at": helper.HNTime(tmp.CreatedAt),
			"created_by": helper.HNString(tmp.CreatedBy),
			"deleted_at": helper.HNTimeGDeletedAt(tmp.DeletedAt),
			"deleted_by": helper.HNString(tmp.DeletedBy),
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}

func (repo *repository) SaveActivityWebhook(activityWebhook ActivityWebhook) (ActivityWebhook, error) {
	if err := repo.DB.Create(&activityWebhook).Error; err != nil {
		return activityWebhook, err
	}
	return activityWebhook, nil
}
