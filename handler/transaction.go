package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/WeAreAmazingTeam/tcd-backend/campaign"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/logs"
	"github.com/WeAreAmazingTeam/tcd-backend/payment"
	"github.com/WeAreAmazingTeam/tcd-backend/transaction"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	transactionSvc transaction.Service
	campaignSvc    campaign.Service
	paymentSvc     payment.Service
	userSvc        user.Service
	logsSvc        logs.Service
}

func NewTransactionHandler(
	transactionService transaction.Service,
	campaignService campaign.Service,
	paymentService payment.Service,
	userService user.Service,
	logsService logs.Service,
) *transactionHandler {
	return &transactionHandler{
		transactionSvc: transactionService,
		campaignSvc:    campaignService,
		paymentSvc:     paymentService,
		userSvc:        userService,
		logsSvc:        logsService,
	}
}

func (handler *transactionHandler) GetAllTransaction(ctx *gin.Context) {
	transactions, err := handler.transactionSvc.GetAllTransaction(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get transactions failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := transaction.FormatMultipleTransactionData(transactions)
	response := helper.APIResponse(http.StatusOK, "Get transactions successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *transactionHandler) GetTransactionByID(ctx *gin.Context) {
	var req transaction.RequestGetTransactionByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get detail transaction failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	transactionDetail, err := handler.transactionSvc.GetTransactionByID(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get detail transaction failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get detail transaction failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := transaction.FormatTransactionData(transactionDetail)
	response := helper.APIResponse(http.StatusOK, "Get detail transaction successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *transactionHandler) CreateTransaction(ctx *gin.Context) {
	var req transaction.RequestCreateTransaction

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Create transaction failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	req.User = ctx.MustGet("userData").(user.User)

	reqCampaign := campaign.RequestGetCampaignByID{}
	reqCampaign.ID = req.CampaignID

	campaign, err := handler.campaignSvc.GetCampaignByID(reqCampaign)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Create transaction failed!", fmt.Sprintf("Campaign with ID %d not found!", req.CampaignID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Create transaction failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if req.UserID != req.User.ID {
		response := helper.APIResponseError(http.StatusBadRequest, "Create transaction failed!", "Bad Request!")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if campaign.Status != "active" {
		response := helper.APIResponseError(http.StatusBadRequest, "Create transaction failed!", "This campaign is not active or already finished!")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if _, err := handler.userSvc.GetUserByID(req.UserID); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Create transaction failed!", fmt.Sprintf("User with ID %d not found!", req.UserID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Create transaction failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	newTransactionData, err := handler.transactionSvc.CreateTransaction(req, campaign.Title)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Create transaction failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := transaction.FormatTransactionData(newTransactionData)
	response := helper.APIResponse(http.StatusCreated, "Create transcation successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v creating transaction id %v.", req.User.Name, newTransactionData.ID))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *transactionHandler) DeleteTransaction(ctx *gin.Context) {
	var reqID transaction.RequestGetTransactionByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete transaction failed", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete transaction.RequestDeleteTransaction

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete transaction failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.transactionSvc.DeleteTransaction(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete transaction failed!", fmt.Sprintf("Transaction with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete transaction failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete transaction successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting transaction id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *transactionHandler) GetTransactionByCampaignID(ctx *gin.Context) {
	var req transaction.RequestGetTransactionByCampaignID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get transaction by campaign id failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	transactions, err := handler.transactionSvc.GetTransactionByCampaignID(ctx, req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get transaction by campaign id failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get transaction by campaign id failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := transaction.FormatMultipleTransactionWitUsernNameData(transactions)
	response := helper.APIResponse(http.StatusOK, "Get transaction by campaign id successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *transactionHandler) GetTransactionByUserID(ctx *gin.Context) {
	var req transaction.RequestGetTransactionByUserID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get transaction by user id failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	transactions, err := handler.transactionSvc.GetTransactionByUserID(ctx, req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get transaction by user id failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get transaction by user id failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := transaction.FormatMultipleTransactionData(transactions)
	response := helper.APIResponse(http.StatusOK, "Get transaction by user id successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *transactionHandler) AdminDataTablesTransactions(ctx *gin.Context) {
	dataTablesTransactions, err := handler.transactionSvc.AdminDataTablesTransactions(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables transactions failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesTransactions)
}

func (handler *transactionHandler) TestMidtrans(ctx *gin.Context) {
	paymentUrl, _, _, err := handler.paymentSvc.RequestPayment(payment.Payment{}, user.User{})

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get payment failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(http.StatusOK, "Get payment successfully!", gin.H{
		"payment_url": paymentUrl,
	})

	ctx.JSON(http.StatusOK, response)
}

func (handler *transactionHandler) UserDataTablesTransactions(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(user.User)
	dataTablesTransactions, err := handler.transactionSvc.UserDataTablesTransactions(ctx, userData)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables transactions failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesTransactions)
}

func (handler *transactionHandler) TransactionWebhooks(ctx *gin.Context) {
	var req transaction.MidtransRequest
	rawData := map[string]string{}

	properties, err := ctx.GetRawData()
	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Failed to process midtrans notification!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	handler.logsSvc.CreateActivityWebhook(logs.RequestCreateActivityWebhook{
		Endpoint:      ctx.Request.URL.Path,
		TriggeredFrom: "MIDTRANS",
		Properties:    string(properties),
	})

	json.Unmarshal([]byte(properties), &rawData)

	req.TransactionTime = rawData["transaction_time"]
	req.TransactionStatus = rawData["transaction_status"]
	req.TransactionID = rawData["transaction_id"]
	req.PaymentType = rawData["payment_type"]
	req.OrderID = rawData["order_id"]
	req.GrossAmount = rawData["gross_amount"]
	req.FraudStatus = rawData["fraud_status"]

	err = handler.transactionSvc.ProcessRequestFromMidtrans(req)
	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Failed to process midtrans notification!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, req)
}

func (handler *transactionHandler) CreateAnonymousTransaction(ctx *gin.Context) {
	var req transaction.RequestCreateAnonymousTransaction

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Create transaction failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqCampaign := campaign.RequestGetCampaignByID{}
	reqCampaign.ID = req.CampaignID

	campaign, err := handler.campaignSvc.GetCampaignByID(reqCampaign)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Create transaction failed!", fmt.Sprintf("Campaign with ID %d not found!", req.CampaignID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Create transaction failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if campaign.Status != "active" {
		response := helper.APIResponseError(http.StatusBadRequest, "Create transaction failed!", "This campaign is not active or already finished!")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	req.User = user.User{
		Name:  "Good Person",
		Email: "m.saleh.solahudin@gmail.com",
	}

	newTransactionData, err := handler.transactionSvc.CreateAnonymousTransaction(req, campaign.Title)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Create transaction failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := transaction.FormatTransactionData(newTransactionData)
	response := helper.APIResponse(http.StatusCreated, "Create transcation successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v creating transaction id %v.", req.User.Name, newTransactionData.ID))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *transactionHandler) CreateTransactionWithEMoney(ctx *gin.Context) {
	var req transaction.RequestCreateTransactionWithEMoney

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Donate failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	req.User = ctx.MustGet("userData").(user.User)

	reqCampaign := campaign.RequestGetCampaignByID{}
	reqCampaign.ID = req.CampaignID

	campaign, err := handler.campaignSvc.GetCampaignByID(reqCampaign)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Donate failed!", fmt.Sprintf("Campaign with ID %d not found!", req.CampaignID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Donate failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if req.UserID != req.User.ID {
		response := helper.APIResponseError(http.StatusBadRequest, "Donate failed!", "Bad Request!")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if campaign.Status != "active" {
		response := helper.APIResponseError(http.StatusBadRequest, "Donate failed!", "This campaign is not active or already finished!")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	user, err := handler.userSvc.GetUserByID(req.UserID)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Donate failed!", fmt.Sprintf("User with ID %d not found!", req.UserID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Donate failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if user.EMoney < float64(req.Amount) {
		response := helper.APIResponseError(http.StatusBadRequest, "Donate failed!", "Your e-Money balance is not enough!")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	newTransactionData, err := handler.transactionSvc.CreateTransactionWithEMoney(req, campaign.Title)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Donate failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	{
		templateData := helper.EmailTransactionSuccess{
			CampaignLink: os.Getenv("WEB_URL") + "/donate/" + strconv.Itoa(campaign.ID),
			Name:         user.Name,
			Amount:       helper.FormatRupiah(float64(req.Amount)),
		}
		go helper.SendMail(user.Email, "Thank You For Your Donation!", templateData, "html/transaction_success.html")
	}

	formatData := transaction.FormatTransactionData(newTransactionData)
	response := helper.APIResponse(http.StatusCreated, "Donate successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v creating transaction id %v.", req.User.Name, newTransactionData.ID))

	ctx.JSON(http.StatusCreated, response)
}
