package transaction

type (
	TransactionFormatter struct {
		ID           int    `json:"id"`
		CampaignID   int    `json:"campaign_id"`
		UserID       int    `json:"user_id"`
		Amount       int64  `json:"amount"`
		Status       string `json:"status"`
		Code         string `json:"code"`
		Comment      string `json:"comment"`
		PaymentURL   string `json:"payment_url"`
		PaymentToken string `json:"payment_token"`
	}

	TransactionCampaignFormatter struct {
		ID         int `json:"id"`
		CampaignID int `json:"campaign_id"`
	}

	TransactionWithUserNameFormatter struct {
		TransactionFormatter
		UserName string `json:"user_name"`
	}
)

func FormatTransactionData(transaction Transaction) (response TransactionFormatter) {
	response = TransactionFormatter{
		ID:           transaction.ID,
		CampaignID:   transaction.CampaignID,
		UserID:       transaction.UserID,
		Amount:       transaction.Amount,
		Status:       transaction.Status,
		Code:         transaction.Code,
		Comment:      transaction.Comment,
		PaymentURL:   transaction.PaymentURL,
		PaymentToken: transaction.PaymentToken,
	}

	return response
}

func FormatMultipleTransactionData(transactions []Transaction) (response []TransactionFormatter) {
	for _, val := range transactions {
		tmp := TransactionFormatter{}
		tmp.ID = val.ID
		tmp.CampaignID = val.CampaignID
		tmp.UserID = val.UserID
		tmp.Amount = val.Amount
		tmp.Status = val.Status
		tmp.Code = val.Code
		tmp.Comment = val.Comment
		tmp.PaymentURL = val.PaymentURL
		tmp.PaymentToken = val.PaymentToken

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []TransactionFormatter{}
	}

	return response
}

func FormatMultipleTransactionWitUsernNameData(transactions []TransactionWithUserName) (response []TransactionWithUserNameFormatter) {
	for _, val := range transactions {
		tmp := TransactionWithUserNameFormatter{}
		tmp.ID = val.ID
		tmp.CampaignID = val.CampaignID
		tmp.UserID = val.UserID
		tmp.UserName = val.UserName
		tmp.Amount = val.Amount
		tmp.Status = val.Status
		tmp.Code = val.Code
		tmp.Comment = val.Comment
		tmp.PaymentURL = val.PaymentURL
		tmp.PaymentToken = val.PaymentToken

		response = append(response, tmp)
	}

	if len(response) == 0 {
		return []TransactionWithUserNameFormatter{}
	}

	return response
}
