package handler

import (
	"net/http"

	"github.com/WeAreAmazingTeam/tcd-backend/campaign"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/logs"
	"github.com/WeAreAmazingTeam/tcd-backend/payment"
	"github.com/WeAreAmazingTeam/tcd-backend/transaction"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
)

type webAndCMSHandler struct {
	transactionSvc transaction.Service
	campaignSvc    campaign.Service
	paymentSvc     payment.Service
	userSvc        user.Service
	logsSvc        logs.Service
}

func NewWebAndCMSHandler(
	transactionService transaction.Service,
	campaignService campaign.Service,
	paymentService payment.Service,
	userService user.Service,
	logsService logs.Service,
) *webAndCMSHandler {
	return &webAndCMSHandler{
		transactionSvc: transactionService,
		campaignSvc:    campaignService,
		paymentSvc:     paymentService,
		userSvc:        userService,
		logsSvc:        logsService,
	}
}

func (handler *webAndCMSHandler) GetStatisticsForHomePage(ctx *gin.Context) {
	totalDonation, err := handler.campaignSvc.GetTotalDonation()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for home page failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	donationCompleted, err := handler.campaignSvc.GetDonationCompleted()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for home page failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	totalTransaction, err := handler.transactionSvc.GetTotalTransaction("")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for home page failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	userRegistered, err := handler.userSvc.GetUserRegistered("")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for home page failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(http.StatusOK, "Get statistics for home page successfully!", gin.H{
		"total_donation":     totalDonation,
		"donation_completed": donationCompleted,
		"total_transaction":  totalTransaction,
		"user_registered":    userRegistered,
	})

	ctx.JSON(http.StatusOK, response)
}

func (handler *webAndCMSHandler) GetStatisticsForAdminDashboard(ctx *gin.Context) {
	totalDonation, err := handler.campaignSvc.GetTotalDonation()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	donationCompleted, err := handler.campaignSvc.GetDonationCompleted()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	totalTransaction, err := handler.transactionSvc.GetTotalTransaction("")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	totalTransactionSuccess, err := handler.transactionSvc.GetTotalTransaction("AND status = 'paid'")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	withdrawalRequests, err := handler.userSvc.GetTotalWithdrawalRequest("")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	withdrawalRequestsProcessed, err := handler.userSvc.GetTotalWithdrawalRequest("AND status != 'pending'")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	userRegistered, err := handler.userSvc.GetUserRegistered("")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	userRegisteredRoleAdmin, err := handler.userSvc.GetUserRegistered("AND role = 'admin'")

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get statistics for admin dashboard failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse(http.StatusOK, "Get statistics for admin dashboard successfully!", gin.H{
		"total_donation":                       totalDonation,
		"donation_completed":                   donationCompleted,
		"total_transaction":                    totalTransaction,
		"total_transaction_success":            totalTransactionSuccess,
		"total_withdrawal_requests":            withdrawalRequests,
		"total_withdrawal_requestes_processed": withdrawalRequestsProcessed,
		"user_registered":                      userRegistered,
		"user_admin_registered":                userRegisteredRoleAdmin,
	})

	ctx.JSON(http.StatusOK, response)
}
