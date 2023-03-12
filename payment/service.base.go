package payment

import (
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct {
}

type Service interface {
	// private
	setupGlobalMidtransConfig()
	createSnapTransaction(payment Payment, user user.User) (*snap.Response, string, error)
	snapRequest(payment Payment, user user.User) (*snap.Request, string)

	// public
	RequestPayment(payment Payment, user user.User) (string, string, string, error)
}

func NewService() *service {
	return &service{}
}
