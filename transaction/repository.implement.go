package transaction

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
)

func (repo *repository) GetAllTransaction(ctx *gin.Context) (transactions []Transaction, err error) {
	additionalQuery := ""

	if ctx.Query("search") != "" {
		s := ctx.Query("search")
		additionalQuery += fmt.Sprintf(" AND (code LIKE '%%%v%%' OR payment_token LIKE '%%%v%%')", s, s)
	}

	if ctx.Query("status") != "" {
		additionalQuery += fmt.Sprintf(" AND status = '%v'", ctx.Query("status"))
	}

	if ctx.Query("limit") != "" {
		additionalQuery += fmt.Sprintf(" LIMIT %v", ctx.Query("limit"))

		if ctx.Query("offset") != "" {
			additionalQuery += fmt.Sprintf(" OFFSET %v", ctx.Query("offset"))
		}
	}

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAll + additionalQuery)).Rows()

	if err != nil {
		return transactions, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := Transaction{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.CampaignID,
			&tmp.UserID,
			&tmp.Amount,
			&tmp.Status,
			&tmp.Code,
			&tmp.Comment,
			&tmp.PaymentURL,
			&tmp.PaymentToken,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, tmp)
	}

	return transactions, nil
}

func (repo *repository) GetTransactionByCampaignId(ctx *gin.Context, campaignID int) (transactions []TransactionWithUserName, err error) {
	additionalQuery := ""

	if ctx.Query("search") != "" {
		s := ctx.Query("search")
		additionalQuery += fmt.Sprintf(" AND (code LIKE '%%%v%%' OR payment_token LIKE '%%%v%%')", s, s)
	}

	if ctx.Query("status") != "" {
		additionalQuery += fmt.Sprintf(" AND status = '%v'", ctx.Query("status"))
	}

	if ctx.Query("order_by") != "" && ctx.Query("order_type") != "" {
		additionalQuery += fmt.Sprintf(" ORDER BY %v %v", ctx.Query("order_by"), ctx.Query("order_type"))
	}

	if ctx.Query("limit") != "" {
		additionalQuery += fmt.Sprintf(" LIMIT %v", ctx.Query("limit"))

		if ctx.Query("offset") != "" {
			additionalQuery += fmt.Sprintf(" OFFSET %v", ctx.Query("offset"))
		}
	}

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetTransactionByCampaignId+additionalQuery), campaignID).Rows()

	if err != nil {
		return transactions, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := TransactionWithUserName{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.CampaignID,
			&tmp.UserID,
			&tmp.UserName,
			&tmp.Amount,
			&tmp.Status,
			&tmp.Code,
			&tmp.Comment,
			&tmp.PaymentURL,
			&tmp.PaymentToken,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, tmp)
	}

	return transactions, nil
}

func (repo *repository) GetTransactionByUserID(ctx *gin.Context, userID int) (transactions []Transaction, err error) {
	additionalQuery := ""

	if ctx.Query("search") != "" {
		s := ctx.Query("search")
		additionalQuery += fmt.Sprintf(" AND (code LIKE '%%%v%%' OR payment_token LIKE '%%%v%%')", s, s)
	}

	if ctx.Query("status") != "" {
		additionalQuery += fmt.Sprintf(" AND status = '%v'", ctx.Query("status"))
	}

	if ctx.Query("limit") != "" {
		additionalQuery += fmt.Sprintf(" LIMIT %v", ctx.Query("limit"))

		if ctx.Query("offset") != "" {
			additionalQuery += fmt.Sprintf(" OFFSET %v", ctx.Query("offset"))
		}
	}

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetTransactionByUserId+additionalQuery), userID).Rows()

	if err != nil {
		return transactions, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := Transaction{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.CampaignID,
			&tmp.UserID,
			&tmp.Amount,
			&tmp.Status,
			&tmp.Code,
			&tmp.Comment,
			&tmp.PaymentURL,
			&tmp.PaymentToken,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, tmp)
	}

	return transactions, nil
}

func (repo *repository) GetTransactionByID(id int) (transaction Transaction, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetTransactionByID), id).Row()

	err = row.Scan(
		&transaction.ID,
		&transaction.CampaignID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Code,
		&transaction.Comment,
		&transaction.PaymentURL,
		&transaction.PaymentToken,
		&transaction.CreatedAt,
		&transaction.CreatedBy,
		&transaction.UpdatedAt,
		&transaction.UpdatedBy,
		&transaction.DeletedAt,
		&transaction.DeletedBy,
	)

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (repo *repository) GetTransactionByCode(code string) (transaction Transaction, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetTransactionByCode), code).Row()

	err = row.Scan(
		&transaction.ID,
		&transaction.CampaignID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Status,
		&transaction.Code,
		&transaction.Comment,
		&transaction.PaymentURL,
		&transaction.PaymentToken,
		&transaction.CreatedAt,
		&transaction.CreatedBy,
		&transaction.UpdatedAt,
		&transaction.UpdatedBy,
		&transaction.DeletedAt,
		&transaction.DeletedBy,
	)

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (repo *repository) SaveTransaction(transaction Transaction) (Transaction, error) {
	if err := repo.DB.Create(&transaction).Error; err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (repo *repository) UpdateTransaction(transaction Transaction) (Transaction, error) {
	if err := repo.DB.Save(&transaction).Error; err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (repo *repository) DeleteTransaction(transaction Transaction) (bool, error) {
	if constant.DELETED_BY {
		if err := repo.DB.Save(&transaction).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	if err := repo.DB.Delete(&transaction).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) AdminDataTablesTransactions(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesTransactions
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

	listOrder := []string{"", "COALESCE((SELECT title FROM campaigns WHERE id = campaign_id), '')", "COALESCE((SELECT name FROM users WHERE id = user_id), '')", "amount", "status", "code", "", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf(`(
			campaign_id LIKE '%%%s%%'
			OR user_id LIKE '%%%s%%'
			OR amount LIKE '%%%s%%'
			OR status LIKE '%%%s%%'
			OR code LIKE '%%%s%%'
			OR COALESCE((SELECT title FROM campaigns WHERE id = campaign_id), '') LIKE '%%%s%%'
			OR COALESCE((SELECT name FROM users WHERE id = winner_user_id), '') LIKE '%%%s%%'
		)`, searchValue,
			searchValue,
			searchValue,
			searchValue,
			searchValue,
			searchValue,
			searchValue,
		)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesTransactions), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesTransactions)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesTransactions)).Scan(&filtered).Error; err != nil {
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
		var (
			campaignName string
			userName     string
		)

		tmp := Transaction{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.CampaignID,
			&campaignName,
			&tmp.UserID,
			&userName,
			&tmp.Amount,
			&tmp.Status,
			&tmp.Code,
			&tmp.Comment,
			&tmp.PaymentURL,
			&tmp.PaymentToken,
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
			"no":            no,
			"id":            tmp.ID,
			"campaign_id":   tmp.CampaignID,
			"campaign_name": campaignName,
			"user_id":       tmp.UserID,
			"user_name":     userName,
			"amount":        tmp.Amount,
			"status":        tmp.Status,
			"code":          tmp.Code,
			"comment":       tmp.Comment,
			"payment_url":   tmp.PaymentURL,
			"payment_token": tmp.PaymentToken,
			"created_at":    helper.HNTime(tmp.CreatedAt),
			"created_by":    helper.HNString(tmp.CreatedBy),
			"updated_at":    helper.HNTime(tmp.UpdatedAt),
			"updated_by":    helper.HNString(tmp.UpdatedBy),
			"deleted_at":    helper.HNTimeGDeletedAt(tmp.DeletedAt),
			"deleted_by":    helper.HNString(tmp.DeletedBy),
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}

func (repo *repository) UserDataTablesTransactions(ctx *gin.Context, user user.User) (result helper.DataTables, err error) {
	var (
		query string = QueryUserDataTablesTransactions
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

	listOrder := []string{"", "campaign_id", "amount", "status", "code", "", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf(`(
			campaign_id LIKE '%%%s%%'
			OR amount LIKE '%%%s%%'
			OR status LIKE '%%%s%%'
			OR code LIKE '%%%s%%'
			OR COALESCE((SELECT title FROM campaigns WHERE id = campaign_id), '') LIKE '%%%s%%'
		)`, searchValue,
			searchValue,
			searchValue,
			searchValue,
			searchValue,
		)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllUserDataTablesTransactions), likes), user.ID).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllUserDataTablesTransactions), user.ID).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllUserDataTablesTransactions), user.ID).Scan(&filtered).Error; err != nil {
			return result, err
		}

		total = filtered
	}

	if order != "" {
		query = fmt.Sprintf("%s %s", query, order)
	}

	query = fmt.Sprintf("%s %s", query, limit)

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(query), user.ID).Rows()

	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			campaignName string
		)

		tmp := Transaction{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.CampaignID,
			&campaignName,
			&tmp.UserID,
			&tmp.Amount,
			&tmp.Status,
			&tmp.Code,
			&tmp.Comment,
			&tmp.PaymentURL,
			&tmp.PaymentToken,
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
			"no":            no,
			"id":            tmp.ID,
			"campaign_id":   tmp.CampaignID,
			"campaign_name": campaignName,
			"user_id":       tmp.UserID,
			"amount":        tmp.Amount,
			"status":        tmp.Status,
			"code":          tmp.Code,
			"comment":       tmp.Comment,
			"payment_url":   tmp.PaymentURL,
			"payment_token": tmp.PaymentToken,
			"created_at":    helper.HNTime(tmp.CreatedAt),
			"created_by":    helper.HNString(tmp.CreatedBy),
			"updated_at":    helper.HNTime(tmp.UpdatedAt),
			"updated_by":    helper.HNString(tmp.UpdatedBy),
			"deleted_at":    helper.HNTimeGDeletedAt(tmp.DeletedAt),
			"deleted_by":    helper.HNString(tmp.DeletedBy),
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}

func (repo *repository) GetTotalTransaction(condition string) (res int, err error) {
	if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetTotalTransaction) + condition).Scan(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}
