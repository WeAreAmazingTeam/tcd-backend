package campaign

import (
	"time"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"
)

type (
	CampaignImage struct {
		ID           int    `json:"id"`
		CampaignID   int    `json:"campaign_id"`
		FileLocation string `json:"file_location"`
		IsPrimary    int    `json:"is_primary"`
		constant.CreatedUpdatedDeleted
	}

	CampaignCategory struct {
		ID       int    `json:"id"`
		Category string `json:"category"`
		constant.CreatedUpdatedDeleted
	}
	Campaign struct {
		ID               int       `json:"id"`
		UserID           int       `json:"user_id"`
		CategoryID       int       `json:"category_id"`
		Title            string    `json:"title"`
		Slug             string    `json:"slug"`
		ShortDescription string    `json:"short_description"`
		Description      string    `json:"description"`
		GoalAmount       int64     `json:"goal_amount"`
		CurrentAmount    int64     `json:"current_amount"`
		IsExclusive      int       `json:"is_exclusive"`
		DonorCount       int       `json:"donor_count"`
		Status           string    `json:"status"`
		FinishedAt       time.Time `json:"finished_at"`
		CampaignImages   []CampaignImage
		constant.CreatedUpdatedDeleted
	}

	ExclusiveCampaign struct {
		ID            int    `json:"id"`
		CampaignID    int    `json:"campaign_id"`
		WinnerUserID  int    `json:"winner_user_id"`
		IsRewardMoney int    `json:"is_reward_money"`
		Reward        string `json:"reward"`
		IsPaidOff     int    `json:"is_paid_off"`
		constant.CreatedUpdatedDeleted
	}
)
