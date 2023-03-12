package campaign

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"

	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (repo *repository) GetAllCampaign(ctx *gin.Context) (campaigns []Campaign, err error) {
	additionalQuery := ""

	if ctx.Query("search") != "" {
		s := ctx.Query("search")
		additionalQuery += fmt.Sprintf(" AND (title LIKE '%%%v%%' OR short_description LIKE '%%%v%%')", s, s)
	}

	if ctx.Query("status") != "" {
		additionalQuery += fmt.Sprintf(" AND status = '%v'", ctx.Query("status"))
	}

	if ctx.Query("category") != "" {
		additionalQuery += fmt.Sprintf(" AND category_id = %v", ctx.Query("category"))
	}

	if ctx.Query("order_by") != "" && ctx.Query("order_type") != "" {
		additionalQuery += fmt.Sprintf(" ORDER BY is_exclusive DESC, %v %v", ctx.Query("order_by"), ctx.Query("order_type"))
	} else {
		additionalQuery += " ORDER BY is_exclusive DESC"
	}

	if ctx.Query("limit") != "" {
		additionalQuery += fmt.Sprintf(" LIMIT %v", ctx.Query("limit"))

		if ctx.Query("offset") != "" {
			additionalQuery += fmt.Sprintf(" OFFSET %v", ctx.Query("offset"))
		}
	}

	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAll + additionalQuery)).Rows()

	if err != nil {
		return campaigns, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := Campaign{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.UserID,
			&tmp.CategoryID,
			&tmp.Title,
			&tmp.Slug,
			&tmp.ShortDescription,
			&tmp.Description,
			&tmp.GoalAmount,
			&tmp.CurrentAmount,
			&tmp.IsExclusive,
			&tmp.DonorCount,
			&tmp.Status,
			&tmp.FinishedAt,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return campaigns, err
		}

		imageRows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetCampaignImages), tmp.ID).Rows()

		if err != nil {
			return campaigns, err
		}

		defer imageRows.Close()

		campaignImages := []CampaignImage{}

		for imageRows.Next() {
			tmpCampaignImages := CampaignImage{}
			errCampaignImages := imageRows.Scan(
				&tmpCampaignImages.ID,
				&tmpCampaignImages.CampaignID,
				&tmpCampaignImages.FileLocation,
				&tmpCampaignImages.IsPrimary,
				&tmpCampaignImages.CreatedAt,
				&tmpCampaignImages.CreatedBy,
				&tmpCampaignImages.UpdatedAt,
				&tmpCampaignImages.UpdatedBy,
				&tmpCampaignImages.DeletedAt,
				&tmpCampaignImages.DeletedBy,
			)

			if errCampaignImages != nil {
				return campaigns, err
			}

			campaignImages = append(campaignImages, tmpCampaignImages)
		}

		tmp.CampaignImages = campaignImages
		campaigns = append(campaigns, tmp)
	}

	return campaigns, nil
}

func (repo *repository) GetCampaignByID(id int) (campaign Campaign, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetCampaignByID), id).Row()

	err = row.Scan(
		&campaign.ID,
		&campaign.UserID,
		&campaign.CategoryID,
		&campaign.Title,
		&campaign.Slug,
		&campaign.ShortDescription,
		&campaign.Description,
		&campaign.GoalAmount,
		&campaign.CurrentAmount,
		&campaign.IsExclusive,
		&campaign.DonorCount,
		&campaign.Status,
		&campaign.FinishedAt,
		&campaign.CreatedAt,
		&campaign.CreatedBy,
		&campaign.UpdatedAt,
		&campaign.UpdatedBy,
		&campaign.DeletedAt,
		&campaign.DeletedBy,
	)

	if err != nil {
		return campaign, err
	}

	imageRows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetCampaignImages), campaign.ID).Rows()

	if err != nil {
		return campaign, err
	}

	defer imageRows.Close()

	campaignImages := []CampaignImage{}

	for imageRows.Next() {
		tmpCampaignImages := CampaignImage{}
		errCampaignImages := imageRows.Scan(
			&tmpCampaignImages.ID,
			&tmpCampaignImages.CampaignID,
			&tmpCampaignImages.FileLocation,
			&tmpCampaignImages.IsPrimary,
			&tmpCampaignImages.CreatedAt,
			&tmpCampaignImages.CreatedBy,
			&tmpCampaignImages.UpdatedAt,
			&tmpCampaignImages.UpdatedBy,
			&tmpCampaignImages.DeletedAt,
			&tmpCampaignImages.DeletedBy,
		)

		if errCampaignImages != nil {
			return campaign, err
		}

		campaignImages = append(campaignImages, tmpCampaignImages)
	}

	campaign.CampaignImages = campaignImages

	return campaign, nil
}

func (repo *repository) SaveCampaign(campaign Campaign) (Campaign, error) {
	if err := repo.DB.Create(&campaign).Error; err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (repo *repository) UpdateCampaign(campaign Campaign) (Campaign, error) {
	if err := repo.DB.Save(&campaign).Error; err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (repo *repository) UpdateCampaignFromPayment(campaignID int, transactionAmount int64) error {
	if err := repo.DB.Model(&Campaign{}).Where("id = ?", campaignID).Updates(map[string]any{"current_amount": gorm.Expr("current_amount + ?", transactionAmount), "donor_count": gorm.Expr("donor_count + ?", 1)}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *repository) DeleteCampaign(campaign Campaign) (bool, error) {
	if constant.DELETED_BY {
		if err := repo.DB.Save(&campaign).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	if err := repo.DB.Delete(&campaign).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) CreateCampaignImage(campaignImage CampaignImage) (CampaignImage, error) {
	err := repo.DB.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}
	return campaignImage, nil
}

func (repo *repository) UpdateAllImagesAsNonPrimary(campaignID int) (bool, error) {
	err := repo.DB.Model(&CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) GetAllCampaignImage() (campaignImages []CampaignImage, err error) {
	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAllImage)).Rows()

	if err != nil {
		return campaignImages, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := CampaignImage{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.CampaignID,
			&tmp.FileLocation,
			&tmp.IsPrimary,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return campaignImages, err
		}

		campaignImages = append(campaignImages, tmp)
	}

	return campaignImages, nil
}

func (repo *repository) GetCampaignImageByID(id int) (campaignImage CampaignImage, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetCampaignImageByID), id).Row()

	err = row.Scan(
		&campaignImage.ID,
		&campaignImage.CampaignID,
		&campaignImage.FileLocation,
		&campaignImage.IsPrimary,
		&campaignImage.CreatedAt,
		&campaignImage.CreatedBy,
		&campaignImage.UpdatedAt,
		&campaignImage.UpdatedBy,
		&campaignImage.DeletedAt,
		&campaignImage.DeletedBy,
	)

	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}

func (repo *repository) DeleteCampaignImage(campaignImage CampaignImage) (bool, error) {
	if err := repo.DB.Delete(&campaignImage).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) GetAllCampaignCategory() (campaignCategories []CampaignCategory, err error) {
	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAllCategory)).Rows()

	if err != nil {
		return campaignCategories, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := CampaignCategory{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Category,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return campaignCategories, err
		}

		campaignCategories = append(campaignCategories, tmp)
	}

	return campaignCategories, nil
}

func (repo *repository) GetCampaignCategoryByID(id int) (campaignCategory CampaignCategory, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetCampaignCategoryByID), id).Row()

	err = row.Scan(
		&campaignCategory.ID,
		&campaignCategory.Category,
		&campaignCategory.CreatedAt,
		&campaignCategory.CreatedBy,
		&campaignCategory.UpdatedAt,
		&campaignCategory.UpdatedBy,
		&campaignCategory.DeletedAt,
		&campaignCategory.DeletedBy,
	)

	if err != nil {
		return campaignCategory, err
	}

	return campaignCategory, nil
}

func (repo *repository) DeleteCampaignCategory(campaignCategory CampaignCategory) (bool, error) {
	tmpCampaignCategory := CampaignCategory{}

	if err := repo.DB.Where("id = ?", campaignCategory.ID).Find(&tmpCampaignCategory).Error; err != nil {
		return false, err
	}

	if tmpCampaignCategory.ID == 0 {
		return false, errors.New("sql: no rows in result set")
	}

	if constant.DELETED_BY {
		if err := repo.DB.Save(&campaignCategory).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	if err := repo.DB.Delete(&campaignCategory).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) SaveCampaignCategory(category CampaignCategory) (CampaignCategory, error) {
	if err := repo.DB.Create(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (repo *repository) UpdateCampaignCategory(category CampaignCategory) (CampaignCategory, error) {
	tmpCampaignCategory := CampaignCategory{}

	if err := repo.DB.Where("id = ?", category.ID).Find(&tmpCampaignCategory).Error; err != nil {
		return category, err
	}

	if tmpCampaignCategory.ID == 0 {
		return category, errors.New("sql: no rows in result set")
	}

	if err := repo.DB.Save(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (repo *repository) GetAllCampaignExclusive() (exclusiveCampaigns []ExclusiveCampaign, err error) {
	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetAllExclusiveCampaign)).Rows()

	if err != nil {
		return exclusiveCampaigns, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := ExclusiveCampaign{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.CampaignID,
			&tmp.WinnerUserID,
			&tmp.IsRewardMoney,
			&tmp.Reward,
			&tmp.IsPaidOff,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return exclusiveCampaigns, err
		}

		exclusiveCampaigns = append(exclusiveCampaigns, tmp)
	}

	return exclusiveCampaigns, nil
}

func (repo *repository) SaveCampaignExclusive(exclusiveCampaign ExclusiveCampaign) (ExclusiveCampaign, error) {
	if err := repo.DB.Create(&exclusiveCampaign).Error; err != nil {
		return exclusiveCampaign, err
	}
	return exclusiveCampaign, nil
}

func (repo *repository) GetCampaignExclusiveByID(id int) (exclusiveCampaign ExclusiveCampaign, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetCampaignExclusiveByID), id).Row()

	err = row.Scan(
		&exclusiveCampaign.ID,
		&exclusiveCampaign.CampaignID,
		&exclusiveCampaign.WinnerUserID,
		&exclusiveCampaign.IsRewardMoney,
		&exclusiveCampaign.Reward,
		&exclusiveCampaign.IsPaidOff,
		&exclusiveCampaign.CreatedAt,
		&exclusiveCampaign.CreatedBy,
		&exclusiveCampaign.UpdatedAt,
		&exclusiveCampaign.UpdatedBy,
		&exclusiveCampaign.DeletedAt,
		&exclusiveCampaign.DeletedBy,
	)

	if err != nil {
		return exclusiveCampaign, err
	}

	return exclusiveCampaign, nil
}

func (repo *repository) GetCampaignExclusiveByCampaignID(id int) (exclusiveCampaign ExclusiveCampaign, err error) {
	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetCampaignExclusiveByCampaignID), id).Row()

	err = row.Scan(
		&exclusiveCampaign.ID,
		&exclusiveCampaign.CampaignID,
		&exclusiveCampaign.WinnerUserID,
		&exclusiveCampaign.IsRewardMoney,
		&exclusiveCampaign.Reward,
		&exclusiveCampaign.IsPaidOff,
		&exclusiveCampaign.CreatedAt,
		&exclusiveCampaign.CreatedBy,
		&exclusiveCampaign.UpdatedAt,
		&exclusiveCampaign.UpdatedBy,
		&exclusiveCampaign.DeletedAt,
		&exclusiveCampaign.DeletedBy,
	)

	if err != nil {
		return exclusiveCampaign, err
	}

	return exclusiveCampaign, nil
}

func (repo *repository) GetCampaignExclusiveByWinnerUserID(id int) (exclusiveCampaigns []ExclusiveCampaign, err error) {
	rows, err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetCampaignExclusiveByWinnerUserID), id).Rows()

	if err != nil {
		return exclusiveCampaigns, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := ExclusiveCampaign{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.CampaignID,
			&tmp.WinnerUserID,
			&tmp.IsRewardMoney,
			&tmp.Reward,
			&tmp.IsPaidOff,
			&tmp.CreatedAt,
			&tmp.CreatedBy,
			&tmp.UpdatedAt,
			&tmp.UpdatedBy,
			&tmp.DeletedAt,
			&tmp.DeletedBy,
		)

		if err != nil {
			return exclusiveCampaigns, err
		}

		exclusiveCampaigns = append(exclusiveCampaigns, tmp)
	}

	return exclusiveCampaigns, nil
}

func (repo *repository) UpdateCampaignExclusive(exclusiveCampaign ExclusiveCampaign) (ExclusiveCampaign, error) {
	if err := repo.DB.Save(&exclusiveCampaign).Error; err != nil {
		return exclusiveCampaign, err
	}
	return exclusiveCampaign, nil
}

func (repo *repository) DeleteCampaignExclusive(exclusiveCampaign ExclusiveCampaign) (bool, error) {
	if err := repo.DB.Delete(&exclusiveCampaign).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (repo *repository) AdminDataTablesCampaigns(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesCampaigns
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

	listOrder := []string{"", "title", "goal_amount", "current_amount", "total_image", "is_exclusive", "status", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(title LIKE '%%%s%%' OR short_description LIKE '%%%s%%' OR status LIKE '%%%s%%')", searchValue, searchValue, searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCampaigns), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCampaigns)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCampaigns)).Scan(&filtered).Error; err != nil {
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
			totalImage int = 0
			finishedAt sql.NullTime
			status     string
		)

		tmp := Campaign{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.UserID,
			&tmp.CategoryID,
			&tmp.Title,
			&tmp.Slug,
			&tmp.ShortDescription,
			&tmp.Description,
			&tmp.GoalAmount,
			&tmp.CurrentAmount,
			&tmp.IsExclusive,
			&tmp.DonorCount,
			&totalImage,
			&status,
			&finishedAt,
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
			"no":                no,
			"id":                tmp.ID,
			"user_id":           tmp.UserID,
			"category_id":       tmp.CategoryID,
			"title":             tmp.Title,
			"slug":              tmp.Slug,
			"short_description": tmp.ShortDescription,
			"description":       tmp.Description,
			"goal_amount":       tmp.GoalAmount,
			"current_amount":    tmp.CurrentAmount,
			"is_exclusive":      tmp.IsExclusive,
			"donor_count":       tmp.DonorCount,
			"total_image":       totalImage,
			"status":            status,
			"finished_at":       helper.HNTime(finishedAt),
			"created_at":        helper.HNTime(tmp.CreatedAt),
			"created_by":        helper.HNString(tmp.CreatedBy),
			"updated_at":        helper.HNTime(tmp.UpdatedAt),
			"updated_by":        helper.HNString(tmp.UpdatedBy),
			"deleted_at":        helper.HNTimeGDeletedAt(tmp.DeletedAt),
			"deleted_by":        helper.HNString(tmp.DeletedBy),
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}

func (repo *repository) AdminDataTablesCategories(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesCategories
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

	listOrder := []string{"", "category", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(category LIKE '%%%s%%')", searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCategories), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCategories)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesCategories)).Scan(&filtered).Error; err != nil {
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
		tmp := CampaignCategory{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.Category,
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
			"category":   tmp.Category,
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

func (repo *repository) UserDataTablesCampaigns(ctx *gin.Context, user user.User) (result helper.DataTables, err error) {
	var (
		query string = QueryUserDataTablesCampaigns
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

	listOrder := []string{"", "title", "goal_amount", "current_amount", "total_image", "is_exclusive", "status", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(title LIKE '%%%s%%' OR short_description LIKE '%%%s%%' OR status LIKE '%%%s%%')", searchValue, searchValue, searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllUserDataTablesCampaigns), likes), user.ID).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllUserDataTablesCampaigns), user.ID).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllUserDataTablesCampaigns), user.ID).Scan(&filtered).Error; err != nil {
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
			totalImage int = 0
			finishedAt sql.NullTime
			status     string
		)

		tmp := Campaign{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.UserID,
			&tmp.CategoryID,
			&tmp.Title,
			&tmp.Slug,
			&tmp.ShortDescription,
			&tmp.Description,
			&tmp.GoalAmount,
			&tmp.CurrentAmount,
			&tmp.IsExclusive,
			&tmp.DonorCount,
			&totalImage,
			&status,
			&finishedAt,
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
			"no":                no,
			"id":                tmp.ID,
			"user_id":           tmp.UserID,
			"category_id":       tmp.CategoryID,
			"title":             tmp.Title,
			"slug":              tmp.Slug,
			"short_description": tmp.ShortDescription,
			"description":       tmp.Description,
			"goal_amount":       tmp.GoalAmount,
			"current_amount":    tmp.CurrentAmount,
			"is_exclusive":      tmp.IsExclusive,
			"donor_count":       tmp.DonorCount,
			"total_image":       totalImage,
			"status":            status,
			"finished_at":       helper.HNTime(finishedAt),
			"created_at":        helper.HNTime(tmp.CreatedAt),
			"created_by":        helper.HNString(tmp.CreatedBy),
			"updated_at":        helper.HNTime(tmp.UpdatedAt),
			"updated_by":        helper.HNString(tmp.UpdatedBy),
			"deleted_at":        helper.HNTimeGDeletedAt(tmp.DeletedAt),
			"deleted_by":        helper.HNString(tmp.DeletedBy),
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}

func (repo *repository) GetTotalDonation() (res int, err error) {
	if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetTotalDonation)).Scan(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (repo *repository) GetDonationCompleted() (res int, err error) {
	if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetDonationCompleted)).Scan(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (repo *repository) AdminDataTablesWinnersExclusiveCampaigns(ctx *gin.Context) (result helper.DataTables, err error) {
	var (
		query string = QueryAdminDataTablesWinnersExclusiveCampaigns
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

	listOrder := []string{"", "COALESCE((SELECT title FROM campaigns WHERE id = campaign_id), '')", "COALESCE((SELECT name FROM users WHERE id = winner_user_id), '')", "reward", "is_paid_off", ""}

	searchValue := ctx.Query("search[value]")
	orderColumn := ctx.Query("order[0][column]")
	starting, _ := strconv.Atoi(ctx.Query("start"))

	if searchValue != "" {
		likes = fmt.Sprintf("(campaign_id = '%s' OR winner_user_id = '%s' OR COALESCE((SELECT title FROM campaigns WHERE id = campaign_id), '') LIKE '%%%s%%' OR COALESCE((SELECT name FROM users WHERE id = winner_user_id), '') LIKE '%%%s%%' OR reward LIKE '%%%s%%')", searchValue, searchValue, searchValue, searchValue, searchValue)
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

		if err := repo.DB.Raw(fmt.Sprintf("%s AND %s", helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesWinnersExclusiveCampaigns), likes)).Scan(&filtered).Error; err != nil {
			return result, err
		}

		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesWinnersExclusiveCampaigns)).Scan(&total).Error; err != nil {
			return result, err
		}
	} else {
		if err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryCountAllAdminDataTablesWinnersExclusiveCampaigns)).Scan(&filtered).Error; err != nil {
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
			campaignName   string
			userWinnerName string
		)

		tmp := ExclusiveCampaign{}
		err := rows.Scan(
			&tmp.ID,
			&tmp.CampaignID,
			&campaignName,
			&tmp.WinnerUserID,
			&userWinnerName,
			&tmp.IsRewardMoney,
			&tmp.Reward,
			&tmp.IsPaidOff,
		)

		if err != nil {
			return result, err
		}

		data = append(data, map[string]any{
			"no":               no,
			"id":               tmp.ID,
			"campaign_id":      tmp.CampaignID,
			"campaign_name":    campaignName,
			"winner_user_id":   tmp.WinnerUserID,
			"winner_user_name": userWinnerName,
			"is_reward_money":  tmp.IsRewardMoney,
			"reward":           tmp.Reward,
			"is_paid_off":      tmp.IsPaidOff,
		})

		no++
	}

	return helper.BuildDatatTables(data, filtered, total), nil
}

func (repo *repository) GetWinnerCampaignExclusive(exclusiveCampaign ExclusiveCampaign) (winnerUserID int, err error) {
	if exclusiveCampaign.IsPaidOff == 0 {
		err := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetOneRandomUserIDForWinnerExclusiveCampaign), exclusiveCampaign.CampaignID).Row().Scan(&winnerUserID)

		if err != nil {
			return 0, err
		}

		return winnerUserID, nil
	}

	return 0, nil
}
