package transaction

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/WeAreAmazingTeam/tcd-backend/campaign"
	"github.com/WeAreAmazingTeam/tcd-backend/company"
	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/payment"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go/example"
)

func (svc *service) GetAllTransaction(ctx *gin.Context) ([]Transaction, error) {
	transactions, err := svc.repo.GetAllTransaction(ctx)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (svc *service) GetTransactionByID(req RequestGetTransactionByID) (Transaction, error) {
	transaction, err := svc.repo.GetTransactionByID(req.ID)

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (svc *service) GetTransactionByCampaignID(ctx *gin.Context, req RequestGetTransactionByCampaignID) ([]TransactionWithUserName, error) {
	transactions, err := svc.repo.GetTransactionByCampaignId(ctx, req.ID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (svc *service) GetTransactionByUserID(ctx *gin.Context, req RequestGetTransactionByUserID) ([]Transaction, error) {
	transactions, err := svc.repo.GetTransactionByUserID(ctx, req.ID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (svc *service) CreateTransaction(req RequestCreateTransaction, campaignName string) (Transaction, error) {
	transaction := Transaction{}

	transaction.CampaignID = req.CampaignID
	transaction.UserID = req.UserID
	transaction.Amount = req.Amount
	transaction.Comment = req.Comment
	transaction.Status = "pending"
	transaction.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	payment := payment.Payment{}
	payment.CampaignID = req.CampaignID
	payment.CampaignName = campaignName
	payment.Amount = req.Amount

	paymentURL, paymentToken, code, err := svc.paymentSvc.RequestPayment(payment, req.User)

	if err != nil {
		return Transaction{}, err
	}

	newTransactionData, err := svc.repo.SaveTransaction(transaction)

	if err != nil {
		return newTransactionData, err
	}

	newTransactionData.Code = code
	newTransactionData.PaymentURL = paymentURL
	newTransactionData.PaymentToken = paymentToken
	newTransactionData.UpdatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	newTransactionData, err = svc.repo.UpdateTransaction(newTransactionData)

	if err != nil {
		return newTransactionData, err
	}

	if req.UserID != 0 {
		svc.companyRepo.CreateCompanyCashFlow(company.CompanyCashFlow{
			Status: "in",
			Amount: int64(req.Amount),
			Note:   fmt.Sprintf("Donate to campaign id %v by %v.", req.CampaignID, req.UserID),
		})
	} else {
		svc.companyRepo.CreateCompanyCashFlow(company.CompanyCashFlow{
			Status: "in",
			Amount: int64(req.Amount),
			Note:   fmt.Sprintf("Donate to campaign id %v by anonymous.", req.CampaignID),
		})
	}

	return newTransactionData, nil
}

func (svc *service) AdminDataTablesTransactions(ctx *gin.Context) (helper.DataTables, error) {
	dataTablesTransactions, err := svc.repo.AdminDataTablesTransactions(ctx)

	if err != nil {
		return dataTablesTransactions, err
	}

	return dataTablesTransactions, nil
}

func (svc *service) DeleteTransaction(reqDetail RequestGetTransactionByID, reqDelete RequestDeleteTransaction) (bool, error) {
	if constant.DELETED_BY {
		transaction, err := svc.repo.GetTransactionByID(reqDetail.ID)

		if err != nil {
			return false, err
		}

		transaction.UpdatedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))
		transaction.DeletedAt = *helper.SetNowNT()
		transaction.DeletedBy = helper.SetNS(strconv.Itoa(reqDelete.User.ID))

		status, err := svc.repo.DeleteTransaction(transaction)

		if err != nil {
			return status, err
		}

		return status, nil
	}

	if _, err := svc.repo.GetTransactionByID(reqDetail.ID); err != nil {
		return false, err
	}

	transaction := Transaction{}

	transaction.ID = reqDetail.ID
	status, err := svc.repo.DeleteTransaction(transaction)

	if err != nil {
		return status, err
	}

	return status, nil
}

func (svc *service) UserDataTablesTransactions(ctx *gin.Context, user user.User) (helper.DataTables, error) {
	dataTablesTransactions, err := svc.repo.UserDataTablesTransactions(ctx, user)

	if err != nil {
		return dataTablesTransactions, err
	}

	return dataTablesTransactions, nil
}

func (svc *service) ProcessRequestFromMidtrans(req MidtransRequest) error {
	transaction, err := svc.repo.GetTransactionByCode(req.OrderID)

	if err != nil {
		log.Println("[transaction webhooks] error while get transaction by code, err: ", err.Error())
		return err
	}

	if req.PaymentType == "credit_card" && req.TransactionStatus == "capture" && req.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if req.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if req.TransactionStatus == "deny" || req.TransactionStatus == "expire" || req.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := svc.repo.UpdateTransaction(transaction)

	if err != nil {
		log.Println("[transaction webhooks] error while update transaction, err: ", err.Error())
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaignData, err := svc.campaignRepo.GetCampaignByID(updatedTransaction.CampaignID)

		if err != nil {
			log.Println("[transaction webhooks] error while get campaign by id, err: ", err.Error())
			return err
		}

		if campaignData.CurrentAmount+updatedTransaction.Amount >= campaignData.GoalAmount {
			campaignData.Status = "finished"
			campaignData.UpdatedBy = helper.SetNS("MIDTRANS")

			updatedCampaign, err := svc.campaignRepo.UpdateCampaign(campaignData)

			if err != nil {
				log.Println("[transaction webhooks] error while update campaign, err: ", err.Error())
				return err
			}

			sumAmount := float64(campaignData.CurrentAmount + updatedTransaction.Amount)
			forDeducted := sumAmount - math.Round(float64(sumAmount-(sumAmount*(float64(6)/float64(100)))))
			deductedAmount := sumAmount - forDeducted

			if err := svc.userRepo.GiveEMoneyToUser(updatedCampaign.UserID, int(deductedAmount)); err != nil {
				log.Println("[transaction webhooks] error while update campaign, err: ", err.Error())
				return err
			}

			svc.userRepo.CreateEMoneyFlow(user.UserEMoneyFlow{
				UserID: updatedCampaign.UserID,
				Status: "in",
				Amount: int64(sumAmount),
				Note:   fmt.Sprintf("Funds from the donation campaign: %v.", updatedCampaign.Title),
			})

			svc.userRepo.CreateEMoneyFlow(user.UserEMoneyFlow{
				UserID: updatedCampaign.UserID,
				Status: "out",
				Amount: int64(deductedAmount),
				Note:   fmt.Sprintf("Admin fee for the donation campaign: %v.", updatedCampaign.Title),
			})

			svc.companyRepo.CreateCompanyCashFlow(company.CompanyCashFlow{
				Status: "out",
				Amount: int64(sumAmount),
				Note:   fmt.Sprintf("Disburse funds for donation campaign: %v.", updatedCampaign.Title),
			})

			svc.companyRepo.CreateCompanyCashFlow(company.CompanyCashFlow{
				Status: "in",
				Amount: int64(deductedAmount),
				Note:   fmt.Sprintf("Admin fee from donation campaign: %v.", updatedCampaign.Title),
			})

			userOwnerCampaign, err := svc.userRepo.GetUserByID(updatedCampaign.UserID)

			if err != nil {
				return err
			}

			{
				templateData := helper.EmailCampaignFinished{
					Campaign:       updatedCampaign,
					Name:           userOwnerCampaign.Name,
					GoalAmount:     helper.FormatRupiah(float64(updatedCampaign.GoalAmount)),
					CollectedFunds: helper.FormatRupiah(sumAmount),
					AdminFee:       helper.FormatRupiah(forDeducted),
					FinalAmount:    helper.FormatRupiah(deductedAmount),
				}
				go helper.SendMail(userOwnerCampaign.Email, fmt.Sprintf("Your Campaign (%v) Has Finished", updatedCampaign.Title), templateData, "html/campaign_finished.html")
			}

			if campaignData.IsExclusive == 1 {
				var reqCheckAndSetWinnerCampaignExclusive campaign.RequestGetCampaignExclusiveByCampaignID
				reqCheckAndSetWinnerCampaignExclusive.ID = campaignData.ID

				if _, err := svc.campaignSvc.CheckAndSetWinnerCampaignExclusive(reqCheckAndSetWinnerCampaignExclusive); err != nil {
					log.Println("[transaction webhooks] error while check and set winner campaign exclusive, err: ", err.Error())
					return err
				}
			}
		}

		if err := svc.campaignRepo.UpdateCampaignFromPayment(campaignData.ID, updatedTransaction.Amount); err != nil {
			log.Println("[transaction webhooks] error while update campaign from payment, err: ", err.Error())
			return err
		}

		if updatedTransaction.UserID != 0 {
			userTransaction, err := svc.userRepo.GetUserByID(updatedTransaction.UserID)

			if err != nil {
				return err
			}

			{
				templateData := helper.EmailTransactionSuccess{
					CampaignLink: os.Getenv("WEB_URL") + "/donate/" + strconv.Itoa(campaignData.ID),
					Name:         userTransaction.Name,
					Amount:       helper.FormatRupiah(float64(transaction.Amount)),
				}
				go helper.SendMail(userTransaction.Email, "Thank You For Your Donation!", templateData, "html/transaction_success.html")
			}
		}
	}

	return nil
}

func (svc *service) GetTotalTransaction(condition string) (res int, err error) {
	res, err = svc.repo.GetTotalTransaction(condition)

	if err != nil {
		return res, err
	}

	return res, nil
}

func (svc *service) CreateAnonymousTransaction(req RequestCreateAnonymousTransaction, campaignName string) (Transaction, error) {
	transaction := Transaction{}

	transaction.CampaignID = req.CampaignID
	transaction.Amount = req.Amount
	transaction.Comment = req.Comment
	transaction.Status = "pending"
	transaction.CreatedBy = helper.SetNS("Anonymous")

	payment := payment.Payment{}
	payment.CampaignID = req.CampaignID
	payment.CampaignName = campaignName
	payment.Amount = req.Amount

	paymentURL, paymentToken, code, err := svc.paymentSvc.RequestPayment(payment, req.User)

	if err != nil {
		return Transaction{}, err
	}

	newTransactionData, err := svc.repo.SaveTransaction(transaction)

	if err != nil {
		return newTransactionData, err
	}

	newTransactionData.Code = code
	newTransactionData.PaymentURL = paymentURL
	newTransactionData.PaymentToken = paymentToken
	newTransactionData.UpdatedBy = helper.SetNS(strconv.Itoa(req.User.ID))
	newTransactionData.UserID = -1

	newTransactionData, err = svc.repo.UpdateTransaction(newTransactionData)

	if err != nil {
		return newTransactionData, err
	}

	return newTransactionData, nil
}

func (svc *service) CreateTransactionWithEMoney(req RequestCreateTransactionWithEMoney, campaignName string) (Transaction, error) {
	transaction := Transaction{}

	transaction.CampaignID = req.CampaignID
	transaction.UserID = req.UserID
	transaction.Amount = req.Amount
	transaction.Comment = req.Comment
	transaction.Status = "paid"
	transaction.Code = fmt.Sprintf("TCD-EMONEY-%v-%v", req.CampaignID, example.Random())
	transaction.PaymentURL = "-"
	transaction.PaymentToken = "-"
	transaction.CreatedBy = helper.SetNS(strconv.Itoa(req.User.ID))

	newTransactionData, err := svc.repo.SaveTransaction(transaction)

	if err != nil {
		return newTransactionData, err
	}

	campaignData, err := svc.campaignRepo.GetCampaignByID(req.CampaignID)

	if err != nil {
		log.Println("[transaction with e-money] error while get campaign by id, err: ", err.Error())
		return newTransactionData, err
	}

	if campaignData.CurrentAmount+req.Amount >= campaignData.GoalAmount {
		campaignData.Status = "finished"
		campaignData.UpdatedBy = helper.SetNS(fmt.Sprintf("Transaction with e-money by user id %v.", req.User.ID))

		updatedCampaign, err := svc.campaignRepo.UpdateCampaign(campaignData)

		if err != nil {
			log.Println("[transaction with e-money] error while update campaign, err: ", err.Error())
			return newTransactionData, err
		}

		sumAmount := float64(campaignData.CurrentAmount + req.Amount)
		forDeducted := sumAmount - math.Round(float64(sumAmount-(sumAmount*(float64(6)/float64(100)))))
		deductedAmount := sumAmount - forDeducted

		if err := svc.userRepo.GiveEMoneyToUser(updatedCampaign.UserID, int(deductedAmount)); err != nil {
			log.Println("[transaction with e-money] error while update campaign, err: ", err.Error())
			return newTransactionData, err
		}

		svc.userRepo.CreateEMoneyFlow(user.UserEMoneyFlow{
			UserID: updatedCampaign.UserID,
			Status: "in",
			Amount: int64(sumAmount),
			Note:   fmt.Sprintf("Funds from the donation campaign: %v.", updatedCampaign.Title),
		})

		svc.userRepo.CreateEMoneyFlow(user.UserEMoneyFlow{
			UserID: updatedCampaign.UserID,
			Status: "out",
			Amount: int64(deductedAmount),
			Note:   fmt.Sprintf("Admin fee for the donation campaign: %v.", updatedCampaign.Title),
		})

		svc.companyRepo.CreateCompanyCashFlow(company.CompanyCashFlow{
			Status: "out",
			Amount: int64(sumAmount),
			Note:   fmt.Sprintf("Disburse funds for exclusive campaign: %v.", updatedCampaign.Title),
		})

		svc.companyRepo.CreateCompanyCashFlow(company.CompanyCashFlow{
			Status: "in",
			Amount: int64(deductedAmount),
			Note:   fmt.Sprintf("Admin fee from donation campaign: %v.", updatedCampaign.Title),
		})

		userOwnerCampaign, err := svc.userRepo.GetUserByID(updatedCampaign.UserID)

		if err != nil {
			return newTransactionData, err
		}

		{
			templateData := helper.EmailCampaignFinished{
				Campaign:       updatedCampaign,
				Name:           userOwnerCampaign.Name,
				GoalAmount:     helper.FormatRupiah(float64(updatedCampaign.GoalAmount)),
				CollectedFunds: helper.FormatRupiah(sumAmount),
				AdminFee:       helper.FormatRupiah(forDeducted),
				FinalAmount:    helper.FormatRupiah(deductedAmount),
			}
			go helper.SendMail(userOwnerCampaign.Email, fmt.Sprintf("Your Campaign (%v) Has Finished", updatedCampaign.Title), templateData, "html/campaign_finished.html")
		}

		if campaignData.IsExclusive == 1 {
			var reqCheckAndSetWinnerCampaignExclusive campaign.RequestGetCampaignExclusiveByCampaignID
			reqCheckAndSetWinnerCampaignExclusive.ID = campaignData.ID

			if _, err := svc.campaignSvc.CheckAndSetWinnerCampaignExclusive(reqCheckAndSetWinnerCampaignExclusive); err != nil {
				log.Println("[transaction with e-money] error while check and set winner campaign exclusive, err: ", err.Error())
				return newTransactionData, err
			}
		}
	}

	if err := svc.campaignRepo.UpdateCampaignFromPayment(campaignData.ID, req.Amount); err != nil {
		log.Println("[transaction with e-money] error while update campaign from payment, err: ", err.Error())
		return newTransactionData, err
	}

	userData, err := svc.userRepo.GetUserByID(req.User.ID)

	if err != nil {
		return newTransactionData, err
	}

	userData.EMoney = userData.EMoney - float64(req.Amount)
	userData, err = svc.userRepo.UpdateUser(userData)

	if err != nil {
		return newTransactionData, err
	}

	svc.userRepo.CreateEMoneyFlow(user.UserEMoneyFlow{
		UserID: req.UserID,
		Status: "out",
		Amount: int64(req.Amount),
		Note:   fmt.Sprintf("Donate with e-money to campaign id %v.", req.CampaignID),
	})

	svc.companyRepo.CreateCompanyCashFlow(company.CompanyCashFlow{
		Status: "in",
		Amount: int64(req.Amount),
		Note:   fmt.Sprintf("Donate with e-money to campaign id %v by %v.", req.CampaignID, req.UserID),
	})

	return newTransactionData, nil
}
