package transaction

import (
	"github.com/WeAreAmazingTeam/tcd-backend/constant"
)

type (
	Transaction struct {
		ID           int    `json:"id"`
		CampaignID   int    `json:"campaign_id"`
		UserID       int    `json:"user_id" gorm:"default:null"`
		Amount       int64  `json:"amount"`
		Status       string `json:"string"`
		Code         string `json:"code"`
		Comment      string `json:"comment"`
		PaymentURL   string `json:"payment_url"`
		PaymentToken string `json:"payment_token"`
		constant.CreatedUpdatedDeleted
	}

	TransactionWithUserName struct {
		Transaction
		UserName string `json:"user_name"`
	}
)
