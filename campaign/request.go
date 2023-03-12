package campaign

import (
	"github.com/WeAreAmazingTeam/tcd-backend/user"
)

type (
	RequestCreateCampaign struct {
		UserID           int    `json:"user_id"`
		CategoryID       int    `json:"category_id" binding:"required"`
		Title            string `json:"title" binding:"required"`
		ShortDescription string `json:"short_description" binding:"required"`
		Description      string `json:"description" binding:"required"`
		GoalAmount       int64  `json:"goal_amount" binding:"required"`
		FinishedAt       string `json:"finished_at"`
		Status           string `json:"status" binding:"required"`
		User             user.User
	}

	RequestUpdateCampaign struct {
		RequestCreateCampaign
	}

	RequestDeleteCampaign struct {
		User user.User
	}

	RequestGetCampaignByID struct {
		ID int `uri:"id" binding:"required"`
	}

	RequestGetCampaignImageByID struct {
		RequestGetCampaignByID
	}

	RequestGetCampaignCategoryByID struct {
		RequestGetCampaignByID
	}

	RequestCreateCampaignImage struct {
		CampaignID int  `form:"campaign_id" binding:"required"`
		IsPrimary  bool `form:"is_primary"`
		User       user.User
	}

	RequestDeleteCampaignImage struct {
		User user.User
	}

	RequestDeleteCampaignCategory struct {
		User user.User
	}

	RequestCreateCampaignCategory struct {
		Category string `json:"category" binding:"required"`
		User     user.User
	}

	RequestUpdateCampaignCategory struct {
		RequestCreateCampaignCategory
	}

	RequestCreateCampaignExclusive struct {
		CampaignID    int    `json:"campaign_id" binding:"required"`
		WinnerUserID  int    `json:"winner_user_id"`
		IsRewardMoney int    `json:"is_reward_money"`
		Reward        string `json:"reward" binding:"required"`
		IsPaidOff     int    `json:"is_paid_off"`
		User          user.User
	}

	RequestGetCampaignExclusiveByID struct {
		RequestGetCampaignByID
	}

	RequestGetCampaignExclusiveByCampaignID struct {
		RequestGetCampaignByID
	}

	RequestGetCampaignExclusiveByWinnerUserID struct {
		RequestGetCampaignByID
	}

	RequestUpdateCampaignExclusive struct {
		RequestCreateCampaignExclusive
	}

	RequestDeleteCampaignExclusive struct {
		User user.User
	}
)
