package payment

import (
	"fmt"
	"os"
	"strings"

	"github.com/WeAreAmazingTeam/tcd-backend/user"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/example"
	"github.com/midtrans/midtrans-go/snap"
)

func (svc *service) setupGlobalMidtransConfig() {
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox
	midtrans.DefaultLoggerLevel = &midtrans.LoggerImplementation{
		LogLevel: midtrans.LogDebug,
	}
}

func (svc *service) createSnapTransaction(payment Payment, user user.User) (*snap.Response, string, error) {
	req, id := svc.snapRequest(payment, user)
	res, err := snap.CreateTransaction(req)

	if err != nil {
		return nil, id, err
	}

	return res, id, nil
}

func (svc *service) snapRequest(payment Payment, user user.User) (*snap.Request, string) {
	id := fmt.Sprintf("TCD-%v-%v", payment.CampaignID, example.Random())

	snapRequest := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  id,
			GrossAmt: int64(payment.Amount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			LName: "",
			Email: user.Email,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    id,
				Price: int64(payment.Amount),
				Qty:   1,
				Name:  strings.ToLower("transaction the cloud donation"),
			},
		},
	}
	return snapRequest, id
}

func (svc *service) RequestPayment(payment Payment, user user.User) (string, string, string, error) {
	svc.setupGlobalMidtransConfig()
	snapTransaction, id, err := svc.createSnapTransaction(payment, user)

	if err != nil {
		return "", "", id, err
	}

	return snapTransaction.RedirectURL, snapTransaction.Token, id, nil
}
