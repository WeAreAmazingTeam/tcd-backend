package campaign

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/WeAreAmazingTeam/tcd-backend/company"
	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func (svc *service) GetAllCampaign(ctx *gin.Context) ([]Campaign, error) {
	campaigns, err := svc.repo.GetAllCampaign(ctx)

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (svc *service) GetCampaignByID(req RequestGetCampaignByID) (Campaign, error) {
	campaign, err := svc.repo.GetCampaignByID(req.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (svc *service) CreateCampaign(req RequestCreateCampaign) (Campaign, error) {
	campaign := Campaign{}
	campaign.UserID = req.User.ID
	campaign.CategoryID = req.CategoryID
	campaign.Title = req.Title
	campaign.ShortDescription = req.ShortDescription
	campaign.Description = req.Description
	campaign.GoalAmount = req.GoalAmount
	campaign.Status = req.Status

	layoutFormat := "2006-01-02"
	finishedAt, _ := time.Parse(layoutFormat, req.FinishedAt)

	campaign.FinishedAt = finishedAt
	slugCandidate := fmt.Sprintf("%s %v%d", req.Title, time.Now().Unix(), req.User.ID)
	campaign.Slug = slug.Make(slugCandidate)
	campaign.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	newCampaignData, err := svc.repo.SaveCampaign(campaign)

	if err != nil {
		return newCampaignData, err
	}

	return newCampaignData, nil
}

func (svc *service) UpdateCampaign(reqDetail RequestGetCampaignByID, reqUpdate RequestUpdateCampaign) (Campaign, error) {
	campaign, err := svc.repo.GetCampaignByID(reqDetail.ID)

	if err != nil {
		return campaign, err
	}

	if campaign.UserID != reqUpdate.User.ID && reqUpdate.User.Role == "user" {
		return campaign, errors.New("not an owner of the campaign")
	}

	if reqUpdate.User.Role == "user" {
		campaign.UserID = reqUpdate.User.ID
	} else {
		campaign.UserID = reqUpdate.UserID
	}

	campaign.CategoryID = reqUpdate.CategoryID
	campaign.Title = reqUpdate.Title
	campaign.ShortDescription = reqUpdate.ShortDescription
	campaign.Description = reqUpdate.Description
	campaign.GoalAmount = reqUpdate.GoalAmount
	campaign.Status = reqUpdate.Status

	if reqUpdate.FinishedAt != "" {
		layoutFormat := "2006-01-02"
		finishedAt, _ := time.Parse(layoutFormat, reqUpdate.FinishedAt)

		campaign.FinishedAt = finishedAt
	}

	campaign.UpdatedBy = helper.SetNS(strconv.Itoa(reqUpdate.User.ID))

	updatedCampaign, err := svc.repo.UpdateCampaign(campaign)

	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (svc *service) DeleteCampaign(reqDetail RequestGetCampaignByID, reqDelete RequestDeleteCampaign) (bool, error) {
	if constant.DELETED_BY {
		campaign, err := svc.repo.GetCampaignByID(reqDetail.ID)

		if err != nil {
			return false, err
		}

		campaign.UpdatedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))
		campaign.DeletedAt = *helper.SetNowNT()
		campaign.DeletedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))

		status, err := svc.repo.DeleteCampaign(campaign)

		if err != nil {
			return status, err
		}

		return status, nil
	}

	campaign, err := svc.repo.GetCampaignByID(reqDetail.ID)

	if err != nil {
		return false, err
	}

	if campaign.UserID != reqDelete.User.ID && reqDelete.User.Role == "user" {
		return false, errors.New("not an owner of the campaign")
	}

	campaign.ID = reqDetail.ID
	status, err := svc.repo.DeleteCampaign(campaign)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) SaveCampaignImage(req RequestCreateCampaignImage, fileLocation string) (campaignImage CampaignImage, err error) {
	campaign, err := svc.repo.GetCampaignByID(req.CampaignID)

	if err != nil {
		return campaignImage, err
	}

	if campaign.UserID != req.User.ID && req.User.Role == "user" {
		return campaignImage, errors.New("not an owner of the campaign")
	}

	isPrimary := 0
	if req.IsPrimary {
		if _, err := svc.repo.UpdateAllImagesAsNonPrimary(req.CampaignID); err != nil {
			return campaignImage, err
		}
		isPrimary = 1
	}

	campaignImage.CampaignID = req.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileLocation = fileLocation
	campaignImage.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	newCampaignImage, err := svc.repo.CreateCampaignImage(campaignImage)

	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}

func (svc *service) GetAllCampaignImage() ([]CampaignImage, error) {
	campaignImages, err := svc.repo.GetAllCampaignImage()

	if err != nil {
		return campaignImages, err
	}

	return campaignImages, nil
}

func (svc *service) GetCampaignImageByID(req RequestGetCampaignImageByID) (CampaignImage, error) {
	campaignImage, err := svc.repo.GetCampaignImageByID(req.ID)

	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}

func (svc *service) DeleteCampaignImage(reqDetail RequestGetCampaignImageByID, reqDelete RequestDeleteCampaignImage) (bool, error) {
	getCampaignImage, err := svc.repo.GetCampaignImageByID(reqDetail.ID)

	if err != nil {
		return false, err
	}

	campaign, err := svc.repo.GetCampaignByID(getCampaignImage.CampaignID)

	if err != nil {
		return false, err
	}

	if campaign.UserID != reqDelete.User.ID && reqDelete.User.Role == "user" {
		return false, errors.New("not an owner of the campaign image")
	}

	campaignImage := CampaignImage{}
	campaignImage.ID = reqDetail.ID
	status, err := svc.repo.DeleteCampaignImage(campaignImage)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) GetAllCampaignCategory() ([]CampaignCategory, error) {
	campaignCategory, err := svc.repo.GetAllCampaignCategory()

	if err != nil {
		return campaignCategory, err
	}

	return campaignCategory, nil
}

func (svc *service) GetCampaignCategoryByID(req RequestGetCampaignCategoryByID) (CampaignCategory, error) {
	campaignCategory, err := svc.repo.GetCampaignCategoryByID(req.ID)

	if err != nil {
		return campaignCategory, err
	}

	return campaignCategory, nil
}

func (svc *service) DeleteCampaignCategory(reqDetail RequestGetCampaignCategoryByID, reqDelete RequestDeleteCampaignCategory) (bool, error) {
	if constant.DELETED_BY {
		campaignCategory, err := svc.repo.GetCampaignCategoryByID(reqDetail.ID)

		if err != nil {
			return false, err
		}

		campaignCategory.UpdatedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))
		campaignCategory.DeletedAt = *helper.SetNowNT()
		campaignCategory.DeletedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))

		status, err := svc.repo.DeleteCampaignCategory(campaignCategory)

		if err != nil {
			return status, err
		}

		return status, nil
	}

	campaignCategory := CampaignCategory{}
	campaignCategory.ID = reqDetail.ID
	status, err := svc.repo.DeleteCampaignCategory(campaignCategory)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) CreateCampaignCategory(req RequestCreateCampaignCategory) (CampaignCategory, error) {
	category := CampaignCategory{}
	category.Category = req.Category
	category.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	newCampaignCategoryData, err := svc.repo.SaveCampaignCategory(category)

	if err != nil {
		return newCampaignCategoryData, err
	}

	return newCampaignCategoryData, nil
}

func (svc *service) UpdateCampaignCategory(reqDetail RequestGetCampaignCategoryByID, reqUpdate RequestUpdateCampaignCategory) (campaignCategory CampaignCategory, err error) {
	campaignCategory, err = svc.repo.GetCampaignCategoryByID(reqDetail.ID)

	if err != nil {
		return campaignCategory, err
	}

	campaignCategory.ID = reqDetail.ID
	campaignCategory.Category = reqUpdate.Category
	campaignCategory.UpdatedBy = helper.SetNS(strconv.Itoa(reqUpdate.User.ID))

	updatedCampaignCategory, err := svc.repo.UpdateCampaignCategory(campaignCategory)

	if err != nil {
		return updatedCampaignCategory, err
	}

	return updatedCampaignCategory, nil
}

func (svc *service) GetAllCampaignExclusive() ([]ExclusiveCampaign, error) {
	exclusiveCampaigns, err := svc.repo.GetAllCampaignExclusive()

	if err != nil {
		return exclusiveCampaigns, err
	}

	return exclusiveCampaigns, nil
}

func (svc *service) CreateCampaignExclusive(req RequestCreateCampaignExclusive) (ExclusiveCampaign, error) {
	exclusiveCampaign := ExclusiveCampaign{}
	exclusiveCampaign.CampaignID = req.CampaignID
	exclusiveCampaign.IsRewardMoney = req.IsRewardMoney
	exclusiveCampaign.Reward = req.Reward
	exclusiveCampaign.IsPaidOff = req.IsPaidOff
	exclusiveCampaign.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	if req.WinnerUserID != 0 {
		exclusiveCampaign.WinnerUserID = req.WinnerUserID
	}

	_, err := svc.repo.GetCampaignByID(req.CampaignID)

	if err != nil {
		return exclusiveCampaign, err
	}

	newCampaignExclusiveData, err := svc.repo.SaveCampaignExclusive(exclusiveCampaign)

	if err != nil {
		return newCampaignExclusiveData, err
	}

	campaign, err := svc.repo.GetCampaignByID(newCampaignExclusiveData.CampaignID)

	if err != nil {
		return newCampaignExclusiveData, err
	}

	campaign.IsExclusive = 1
	campaign.UpdatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	if _, err := svc.repo.UpdateCampaign(campaign); err != nil {
		return newCampaignExclusiveData, err
	}

	return newCampaignExclusiveData, nil
}

func (svc *service) GetCampaignExclusiveByID(req RequestGetCampaignExclusiveByID) (ExclusiveCampaign, error) {
	exclusiveCampaign, err := svc.repo.GetCampaignExclusiveByID(req.ID)

	if err != nil {
		return exclusiveCampaign, err
	}

	return exclusiveCampaign, nil
}

func (svc *service) GetCampaignExclusiveByCampaignID(req RequestGetCampaignExclusiveByCampaignID) (ExclusiveCampaign, error) {
	exclusiveCampaign, err := svc.repo.GetCampaignExclusiveByCampaignID(req.ID)

	if err != nil {
		return exclusiveCampaign, err
	}

	return exclusiveCampaign, nil
}

func (svc *service) GetCampaignExclusiveByWinnerUserID(req RequestGetCampaignExclusiveByWinnerUserID) ([]ExclusiveCampaign, error) {
	exclusiveCampaigns, err := svc.repo.GetCampaignExclusiveByWinnerUserID(req.ID)

	if err != nil {
		return exclusiveCampaigns, err
	}

	return exclusiveCampaigns, nil
}

func (svc *service) UpdateCampaignExclusive(reqDetail RequestGetCampaignExclusiveByID, reqUpdate RequestUpdateCampaignExclusive) (exclusiveCampaign ExclusiveCampaign, err error) {
	exclusiveCampaign, err = svc.repo.GetCampaignExclusiveByID(reqDetail.ID)

	if err != nil {
		return exclusiveCampaign, err
	}

	exclusiveCampaign.ID = reqDetail.ID
	exclusiveCampaign.CampaignID = reqUpdate.CampaignID
	exclusiveCampaign.WinnerUserID = reqUpdate.WinnerUserID
	exclusiveCampaign.IsRewardMoney = reqUpdate.IsRewardMoney
	exclusiveCampaign.Reward = reqUpdate.Reward
	exclusiveCampaign.IsPaidOff = reqUpdate.IsPaidOff
	exclusiveCampaign.UpdatedBy = helper.SetNS(strconv.Itoa(reqUpdate.User.ID))

	updatedExclusiveCampaign, err := svc.repo.UpdateCampaignExclusive(exclusiveCampaign)

	if err != nil {
		return updatedExclusiveCampaign, err
	}

	return updatedExclusiveCampaign, nil
}

func (svc *service) DeleteCampaignExclusive(reqDetail RequestGetCampaignExclusiveByID, reqDelete RequestDeleteCampaignExclusive) (bool, error) {
	exclusiveCampaign := ExclusiveCampaign{}
	exclusiveCampaign.ID = reqDetail.ID

	prevData, err := svc.repo.GetCampaignExclusiveByID(reqDetail.ID)

	if err != nil {
		return false, err
	}

	status, err := svc.repo.DeleteCampaignExclusive(exclusiveCampaign)

	if err != nil {
		return status, err
	}

	campaign, err := svc.repo.GetCampaignByID(prevData.CampaignID)

	if err != nil {
		return false, err
	}

	campaign.IsExclusive = 0
	campaign.UpdatedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))

	if _, err := svc.repo.UpdateCampaign(campaign); err != nil {
		return false, err
	}

	return status, nil
}

func (svc *service) AdminDataTablesCampaigns(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesCampaigns, err := svc.repo.AdminDataTablesCampaigns(ctx)

	if err != nil {
		return dataTablesCampaigns, err
	}

	return dataTablesCampaigns, nil
}

func (svc *service) AdminDataTablesCategories(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesCategories, err := svc.repo.AdminDataTablesCategories(ctx)

	if err != nil {
		return dataTablesCategories, err
	}

	return dataTablesCategories, nil
}

func (svc *service) UserDataTablesCampaigns(ctx *gin.Context, user user.User) (helper.DataTables, error) {
	dataTablesCampaigns, err := svc.repo.UserDataTablesCampaigns(ctx, user)

	if err != nil {
		return dataTablesCampaigns, err
	}

	return dataTablesCampaigns, nil
}

func (svc *service) GetTotalDonation() (res int, err error) {
	res, err = svc.repo.GetTotalDonation()

	if err != nil {
		return res, err
	}

	return res, nil
}

func (svc *service) GetDonationCompleted() (res int, err error) {
	res, err = svc.repo.GetDonationCompleted()

	if err != nil {
		return res, err
	}

	return res, nil
}

func (svc *service) AdminDataTablesWinnersExclusiveCampaigns(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesWinnerExclusiveCampaigns, err := svc.repo.AdminDataTablesWinnersExclusiveCampaigns(ctx)

	if err != nil {
		return dataTablesWinnerExclusiveCampaigns, err
	}

	return dataTablesWinnerExclusiveCampaigns, nil
}

func (svc *service) CheckAndSetWinnerCampaignExclusive(req RequestGetCampaignExclusiveByCampaignID) (exclusiveCampaign ExclusiveCampaign, err error) {
	exclusiveCampaign, err = svc.repo.GetCampaignExclusiveByCampaignID(req.ID)

	if err != nil {
		return exclusiveCampaign, err
	}

	winnerUserID, err := svc.repo.GetWinnerCampaignExclusive(exclusiveCampaign)

	if err != nil {
		return exclusiveCampaign, err
	}

	if winnerUserID == 0 {
		return exclusiveCampaign, errors.New("no user can be the winner")
	}

	exclusiveCampaign.WinnerUserID = winnerUserID

	if exclusiveCampaign.IsRewardMoney == 1 {
		exclusiveCampaign.IsPaidOff = 1
	}

	updatedExclusiveCampaign, err := svc.repo.UpdateCampaignExclusive(exclusiveCampaign)

	if err != nil {
		return updatedExclusiveCampaign, err
	}

	if updatedExclusiveCampaign.IsRewardMoney == 1 {
		userData, err := svc.userRepo.GetUserByID(winnerUserID)

		if err != nil {
			return updatedExclusiveCampaign, err
		}

		moneyReward, err := strconv.Atoi(updatedExclusiveCampaign.Reward)

		if err != nil {
			return updatedExclusiveCampaign, err
		}

		userData.EMoney = userData.EMoney + float64(moneyReward)
		userData, err = svc.userRepo.UpdateUser(userData)

		if err != nil {
			return updatedExclusiveCampaign, err
		}

		svc.userRepo.CreateEMoneyFlow(user.UserEMoneyFlow{
			UserID: userData.ID,
			Status: "in",
			Amount: int64(moneyReward),
			Note:   fmt.Sprintf("Reward from exclusive campaign id %v.", updatedExclusiveCampaign.CampaignID),
		})

		svc.companyRepo.CreateCompanyCashFlow(company.CompanyCashFlow{
			Status: "out",
			Amount: int64(moneyReward),
			Note:   fmt.Sprintf("Reward for exclusive campaign id %v.", updatedExclusiveCampaign.CampaignID),
		})
	}

	winnerUserData, err := svc.userRepo.GetUserByID(winnerUserID)

	if err != nil {
		return updatedExclusiveCampaign, err
	}

	{
		status := "Pending"

		if updatedExclusiveCampaign.IsPaidOff == 1 {
			status = "Paid Off"
		}

		templateData := helper.EmailEarningRewardFromExclusiveCampaign{
			CampaignLink: os.Getenv("WEB_URL") + "/donate/" + strconv.Itoa(updatedExclusiveCampaign.CampaignID),
			Name:         winnerUserData.Name,
			Reward:       updatedExclusiveCampaign.Reward,
			Status:       status,
		}
		go helper.SendMail(winnerUserData.Email, "Congratulations, You Get Rewards From Exclusive Campaign!", templateData, "html/earn_reward.html")
	}

	return exclusiveCampaign, nil
}
