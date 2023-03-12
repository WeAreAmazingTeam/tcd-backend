package user

import (
	"errors"
	"fmt"
	"os/user"
	"strconv"
	"strings"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (repo *repository) SaveUser(user User) (User, error) {
	if err := repo.DB.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (repo *repository) GetUserByEmail(email string) (user User, err error) {
	if err := repo.DB.Where("email = ?", email).Find(&user).Error; err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("sql: no rows in result set")
	}

	return user, nil
}

func (repo *repository) GetUserByID(id int) (user User, err error) {
	if err := repo.DB.Where("id = ?", id).Find(&user).Error; err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("sql: no rows in result set")
	}

	return user, nil
}

func (repo *repository) GetAllUser() (users []User, err error) {
	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAllUser)).Rows()

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := User{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Role,
			&tmp.Name,
			&tmp.Email,
			&tmp.Password,
			&tmp.EMoney,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return users, err
		}

		users = append(users, tmp)
	}

	return users, nil
}

func (repo *repository) UpdateUser(user User) (User, error) {
	if err := repo.DB.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (repo *repository) DeleteUser(user User) (bool, error) {
	if constant.DELETED_BY {
		if err := repo.DB.Save(&user).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	if err := repo.DB.Delete(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) AdminDataTablesUsers(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesUsers
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

	listOrder := []string{"", "name", "email", "role", "e_money", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf(`(
			name LIKE '%%%s%%' OR email LIKE '%%%s%%' OR role LIKE '%%%s%%' OR e_money LIKE '%%%s%%'
		)`, searchValue, searchValue, searchValue, searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesUsers), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesUsers)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesUsers)).Scan(&filtered).Error; err != nil {
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
		tmp := User{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Role,
			&tmp.Name,
			&tmp.Email,
			&tmp.Password,
			&tmp.EMoney,
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
			"role":       tmp.Role,
			"name":       tmp.Name,
			"email":      tmp.Email,
			"password":   tmp.Password,
			"e_money":    tmp.EMoney,
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

func (repo *repository) GetUserRegistered(condition string) (res int, err error) {
	if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryUserRegistered) + condition).Scan(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (repo *repository) GetTotalWithdrawalRequest(condition string) (res int, err error) {
	if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryTotalWithdrawalRequest) + condition).Scan(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (repo *repository) GiveEMoneyToUser(userID, eMoney int) error {
	result := repo.DB.Model(&user.User{}).Where("id = ?", userID).Update("e_money", gorm.Expr("e_money + ?", eMoney))

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *repository) CreateEMoneyFlow(userEMoneyFlow UserEMoneyFlow) (UserEMoneyFlow, error) {
	if err := repo.DB.Create(&userEMoneyFlow).Error; err != nil {
		return userEMoneyFlow, err
	}
	return userEMoneyFlow, nil
}

func (repo *repository) UserDataTablesEMoneyFlow(ctx *gin.Context, user User) (result helper.DataTables, err error) {
	var (
		query string = QueryDataTablesEMoneyFlow
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

	listOrder := []string{"", "status", "amount", "note", "created_at"}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(status LIKE '%%%s%%' OR amount LIKE '%%%s%%' OR note LIKE '%%%s%%')", searchValue, searchValue, searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllDataTablesEMoneyFlow), likes), user.ID).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllDataTablesEMoneyFlow), user.ID).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllDataTablesEMoneyFlow), user.ID).Scan(&filtered).Error; err != nil {
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
		tmp := UserEMoneyFlow{}
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

func (repo *repository) CreateWithdrawalRequest(userWithdrawalRequest UserWithdrawalRequest) (UserWithdrawalRequest, error) {
	if err := repo.DB.Create(&userWithdrawalRequest).Error; err != nil {
		return userWithdrawalRequest, err
	}
	return userWithdrawalRequest, nil
}

func (repo *repository) UserDataTablesWithdrawalRequest(ctx *gin.Context, user User) (result helper.DataTables, err error) {
	var (
		query string = QueryDataTablesWithdrawalRequest
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

	listOrder := []string{"", "status", "amount", "note", "created_at"}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(status LIKE '%%%s%%' OR amount LIKE '%%%s%%' OR note LIKE '%%%s%%')", searchValue, searchValue, searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllDataTablesWithdrawalRequest), likes), user.ID).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllDataTablesWithdrawalRequest), user.ID).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllDataTablesWithdrawalRequest), user.ID).Scan(&filtered).Error; err != nil {
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
		tmp := UserWithdrawalRequest{}
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

func (repo *repository) AdminDataTablesWithdrawalRequest(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesWithdrawalRequest
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

	listOrder := []string{"", "COALESCE((SELECT name FROM users WHERE id = user_id), '')", "amount", "note", "status", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(COALESCE((SELECT name FROM users WHERE id = user_id), '') LIKE '%%%s%%' OR amount LIKE '%%%s%%' OR note LIKE '%%%s%%' OR status LIKE '%%%s%%')", searchValue, searchValue, searchValue, searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesWithdrawalRequest), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesWithdrawalRequest)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesWithdrawalRequest)).Scan(&filtered).Error; err != nil {
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
		var userName string
		tmp := UserWithdrawalRequest{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.UserID,
			&userName,
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
			"user_id":    tmp.UserID,
			"user_name":  userName,
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

func (repo *repository) GetWithdrawalRequestByID(id int) (userWithdrawalRequest UserWithdrawalRequest, err error) {
	if err := repo.DB.Where("id = ?", id).Find(&userWithdrawalRequest).Error; err != nil {
		return userWithdrawalRequest, err
	}

	if userWithdrawalRequest.ID == 0 {
		return userWithdrawalRequest, errors.New("sql: no rows in result set")
	}

	return userWithdrawalRequest, nil
}

func (repo *repository) UpdateUserWithdrawalRequest(userWithdrawalRequest UserWithdrawalRequest) (UserWithdrawalRequest, error) {
	if err := repo.DB.Save(&userWithdrawalRequest).Error; err != nil {
		return userWithdrawalRequest, err
	}
	return userWithdrawalRequest, nil
}

func (repo *repository) DeleteUserWithdrawalRequest(userWithdrawalRequest UserWithdrawalRequest) (bool, error) {
	if constant.DELETED_BY {
		if err := repo.DB.Save(&userWithdrawalRequest).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	if err := repo.DB.Delete(&userWithdrawalRequest).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) CreateForgotPasswordToken(userForgotPasswordToken UserForgotPasswordToken) (UserForgotPasswordToken, error) {
	if err := repo.DB.Create(&userForgotPasswordToken).Error; err != nil {
		return userForgotPasswordToken, err
	}
	return userForgotPasswordToken, nil
}

func (repo *repository) GetDataForgotPasswordByToken(token string) (userForgotPasswordToken UserForgotPasswordToken, err error) {
	if err := repo.DB.Where("token = ?", token).Find(&userForgotPasswordToken).Error; err != nil {
		return userForgotPasswordToken, err
	}
	return userForgotPasswordToken, nil
}

func (repo *repository) DeleteForgotPasswordToken(userForgotPasswordToken UserForgotPasswordToken) (bool, error) {
	if err := repo.DB.Delete(&userForgotPasswordToken).Error; err != nil {
		return false, err
	}
	return true, nil
}
