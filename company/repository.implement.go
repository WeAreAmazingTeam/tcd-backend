package company

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
)

func (repo *repository) GetCompanyCashFlowByID(id int) (companyCashFlow CompanyCashFlow, err error) {
	if err := repo.DB.Where("id = ?", id).Find(&companyCashFlow).Error; err != nil {
		return companyCashFlow, err
	}
	return companyCashFlow, nil
}

func (repo *repository) CreateCompanyCashFlow(companCashFlow CompanyCashFlow) (CompanyCashFlow, error) {
	if err := repo.DB.Create(&companCashFlow).Error; err != nil {
		return companCashFlow, err
	}
	return companCashFlow, nil
}

func (repo *repository) DeleteCompanyCashFlow(companyCashFlow CompanyCashFlow) (bool, error) {
	tmpCompanyCashFlow := CompanyCashFlow{}

	if err := repo.DB.Where("id = ?", companyCashFlow.ID).Find(&tmpCompanyCashFlow).Error; err != nil {
		return false, err
	}

	if tmpCompanyCashFlow.ID == 0 {
		return false, errors.New("sql: no rows in result set")
	}

	if constant.DELETED_BY {
		if err := repo.DB.Save(&companyCashFlow).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	if err := repo.DB.Delete(&companyCashFlow).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) AdminDataTablesCompanyCashFlow(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesCompanyCashFlow
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

	listOrder := []string{"", "amount", "note", "status", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(amount LIKE '%%%s%%' OR note LIKE '%%%s%%' OR status LIKE '%%%s%%')", searchValue, searchValue, searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCompanyCashFlow), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCompanyCashFlow)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCompanyCashFlow)).Scan(&filtered).Error; err != nil {
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
		tmp := CompanyCashFlow{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Status,
			&tmp.Amount,
			&tmp.Note,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return result, err
		}

		data = append(data, map[string]any{
			"no":         no,
			"id":         tmp.ID,
			"status":     tmp.Status,
			"amount":     tmp.Amount,
			"note":       tmp.Note,
			"created_at": helper.HNTime(tmp.CreatedAt),
			"created_by": helper.HNString(tmp.CreatedBy),
			"updated_at": helper.HNTime(tmp.UpdatedAt),
			"updated_by": helper.HNString(tmp.UpdatedBy),
			"deleted_at": helper.HNTimeGDeletedAt(tmp.DeletedAt),
			"deleted_by": helper.HNString(tmp.DeletedBy),
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}
