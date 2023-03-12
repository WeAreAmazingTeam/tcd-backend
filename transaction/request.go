package transaction

import "github.com/WeAreAmazingTeam/tcd-backend/user"

type (
	RequestCreateTransaction struct {
		CampaignID int    `json:"campaign_id"`
		UserID     int    `json:"user_id"`
		Amount     int64  `json:"amount"`
		Comment    string `json:"comment"`
		User       user.User
	}

	RequestGetTransactionByID struct {
		ID int `uri:"id" binding:"required"`
	}

	RequestGetTransactionByCampaignID struct {
		RequestGetTransactionByID
	}

	RequestGetTransactionByUserID struct {
		RequestGetTransactionByID
	}

	RequestDeleteTransaction struct {
		RequestCreateTransaction
	}

	MidtransRequest struct {
		TransactionTime   string `json:"transaction_time"`
		TransactionStatus string `json:"transaction_status"`
		TransactionID     string `json:"transaction_id"`
		PaymentType       string `json:"payment_type"`
		OrderID           string `json:"order_id"`
		GrossAmount       string `json:"gross_amount"`
		FraudStatus       string `json:"fraud_status"`
	}

	RequestCreateAnonymousTransaction struct {
		CampaignID int    `json:"campaign_id"`
		Amount     int64  `json:"amount"`
		Comment    string `json:"comment"`
		User       user.User
	}

	RequestCreateTransactionWithEMoney struct {
		RequestCreateTransaction
	}
)
