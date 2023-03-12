package company

type (
	CompanyCashFlowFormatter struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
		Amount int64  `json:"amount"`
		Note   string `json:"note"`
	}
)

func FormatCompanyCashFlowData(companyCashFlow CompanyCashFlow) (response CompanyCashFlowFormatter) {
	response = CompanyCashFlowFormatter{
		ID:     companyCashFlow.ID,
		Status: companyCashFlow.Status,
		Amount: companyCashFlow.Amount,
		Note:   companyCashFlow.Note,
	}

	return response
}
