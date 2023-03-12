package campaign

import (
	"time"
)

type (
	CampaignFormatter struct {
		ID               int                                       `json:"id"`
		UserID           int                                       `json:"user_id"`
		CategoryID       int                                       `json:"category_id"`
		Title            string                                    `json:"title"`
		Slug             string                                    `json:"slug"`
		ShortDescription string                                    `json:"short_description"`
		Description      string                                    `json:"description"`
		GoalAmount       int64                                     `json:"goal_amount"`
		CurrentAmount    int64                                     `json:"current_amount"`
		IsExclusive      int                                       `json:"is_exclusive"`
		DonorCount       int                                       `json:"donor_count"`
		Status           string                                    `json:"status"`
		FinishedAt       time.Time                                 `json:"finished_at"`
		CampaignImages   []CampaignImageWithoutCampaignIDFormatter `json:"images"`
	}

	CampaignImageFormatter struct {
		ID           int    `json:"id"`
		CampaignID   int    `json:"campaign_id"`
		FileLocation string `json:"file_location"`
		IsPrimary    int    `json:"is_primary"`
	}

	CampaignImageWithoutCampaignIDFormatter struct {
		ID           int    `json:"id"`
		FileLocation string `json:"file_location"`
		IsPrimary    int    `json:"is_primary"`
	}

	CampaignCategoryFormatter struct {
		ID       int    `json:"id"`
		Category string `json:"category"`
	}

	CampaignExclusiveFormatter struct {
		ID            int    `json:"id"`
		CampaignID    int    `json:"campaign_id"`
		WinnerUserID  int    `json:"winner_user_id"`
		IsRewardMoney int    `json:"is_reward_money"`
		Reward        string `json:"reward"`
		IsPaidOff     int    `json:"is_paid_off"`
	}
)

func FormatCampaignData(campaign Campaign) (response CampaignFormatter) {
	response = CampaignFormatter{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		CategoryID:       campaign.CategoryID,
		Title:            campaign.Title,
		Slug:             campaign.Slug,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		IsExclusive:      campaign.IsExclusive,
		DonorCount:       campaign.DonorCount,
		Status:           campaign.Status,
		FinishedAt:       campaign.FinishedAt,
	}

	images := []CampaignImageWithoutCampaignIDFormatter{}
	tmpImages := CampaignImageWithoutCampaignIDFormatter{}

	for _, img := range campaign.CampaignImages {
		tmpImages.ID = img.ID
		tmpImages.FileLocation = img.FileLocation
		tmpImages.IsPrimary = img.IsPrimary

		images = append(images, tmpImages)
	}

	response.CampaignImages = images

	return response
}

func FormatMultipleCampaignData(campaigns []Campaign) (response []CampaignFormatter) {
	tmp := CampaignFormatter{}
	tmpImages := CampaignImageWithoutCampaignIDFormatter{}

	for _, val := range campaigns {
		tmp.ID = val.ID
		tmp.UserID = val.UserID
		tmp.CategoryID = val.CategoryID
		tmp.Title = val.Title
		tmp.Slug = val.Slug
		tmp.ShortDescription = val.ShortDescription
		tmp.Description = val.Description
		tmp.GoalAmount = val.GoalAmount
		tmp.CurrentAmount = val.CurrentAmount
		tmp.IsExclusive = val.IsExclusive
		tmp.DonorCount = val.DonorCount
		tmp.Status = val.Status
		tmp.FinishedAt = val.FinishedAt

		images := []CampaignImageWithoutCampaignIDFormatter{}

		for _, img := range val.CampaignImages {
			tmpImages.ID = img.ID
			tmpImages.FileLocation = img.FileLocation
			tmpImages.IsPrimary = img.IsPrimary

			images = append(images, tmpImages)
		}

		tmp.CampaignImages = images

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []CampaignFormatter{}
	}

	return response
}

func FormatCampaignImageData(campaign CampaignImage) (response CampaignImageFormatter) {
	response = CampaignImageFormatter{
		ID:           campaign.ID,
		CampaignID:   campaign.CampaignID,
		FileLocation: campaign.FileLocation,
		IsPrimary:    campaign.IsPrimary,
	}

	return response
}

func FormatMultipleCampaignImageData(campaignImages []CampaignImage) (response []CampaignImageFormatter) {
	for _, val := range campaignImages {
		tmp := CampaignImageFormatter{}
		tmp.ID = val.ID
		tmp.CampaignID = val.CampaignID
		tmp.FileLocation = val.FileLocation
		tmp.IsPrimary = val.IsPrimary

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []CampaignImageFormatter{}
	}

	return response
}

func FormatCampaignCategoryData(category CampaignCategory) (response CampaignCategoryFormatter) {
	response = CampaignCategoryFormatter{
		ID:       category.ID,
		Category: category.Category,
	}

	return response
}

func FormatMultipleCampaignCategoryData(campaignCategories []CampaignCategory) (response []CampaignCategoryFormatter) {
	for _, val := range campaignCategories {
		tmp := CampaignCategoryFormatter{}
		tmp.ID = val.ID
		tmp.Category = val.Category

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []CampaignCategoryFormatter{}
	}

	return response
}

func FormatCampaignExclusiveData(exclusiveCampaign ExclusiveCampaign) (response CampaignExclusiveFormatter) {
	response = CampaignExclusiveFormatter{
		ID:            exclusiveCampaign.ID,
		CampaignID:    exclusiveCampaign.CampaignID,
		WinnerUserID:  exclusiveCampaign.WinnerUserID,
		IsRewardMoney: exclusiveCampaign.IsRewardMoney,
		Reward:        exclusiveCampaign.Reward,
		IsPaidOff:     exclusiveCampaign.IsPaidOff,
	}

	return response
}

func FormatMultipleCampaignExclusiveData(exclusiveCampaigns []ExclusiveCampaign) (response []CampaignExclusiveFormatter) {
	for _, val := range exclusiveCampaigns {
		tmp := CampaignExclusiveFormatter{}
		tmp.ID = val.ID
		tmp.CampaignID = val.CampaignID
		tmp.WinnerUserID = val.WinnerUserID
		tmp.IsRewardMoney = val.IsRewardMoney
		tmp.Reward = val.Reward
		tmp.IsPaidOff = val.IsPaidOff

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []CampaignExclusiveFormatter{}
	}

	return response
}
