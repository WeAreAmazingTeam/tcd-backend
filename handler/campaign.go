package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/WeAreAmazingTeam/tcd-backend/campaign"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/logs"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type campaignHandler struct {
	campaignSvc campaign.Service
	userSvc     user.Service
	logsSvc     logs.Service
}

func NewCampaignHandler(
	campaignService campaign.Service,
	userService user.Service,
	logsService logs.Service,
) *campaignHandler {
	return &campaignHandler{
		campaignSvc: campaignService,
		userSvc:     userService,
		logsSvc:     logsService,
	}
}

func (handler *campaignHandler) GetAllCampaign(ctx *gin.Context) {
	campaigns, err := handler.campaignSvc.GetAllCampaign(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get campaigns failed!", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatData := campaign.FormatMultipleCampaignData(campaigns)
	response := helper.APIResponse(http.StatusOK, "Get campaigns successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) GetCampaignByID(ctx *gin.Context) {
	var req campaign.RequestGetCampaignByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get detail campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	campaignDetail, err := handler.campaignSvc.GetCampaignByID(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get detail campaign failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get detail campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignData(campaignDetail)
	response := helper.APIResponse(http.StatusOK, "Get detail campaign successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) CreateCampaign(ctx *gin.Context) {
	var req campaign.RequestCreateCampaign

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Create campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userData := ctx.MustGet("userData").(user.User)

	if req.UserID != 0 {
		if userData.Role == "user" {
			response := helper.APIResponseError(http.StatusBadRequest, "Create campaign failed!", "Bad Request!")
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		if _, err := handler.userSvc.GetUserByID(req.UserID); err != nil {
			if helper.IsErrNoRows(err.Error()) {
				response := helper.APIResponseError(http.StatusNotFound, "Create campaign failed!", fmt.Sprintf("User with ID %d not found!", req.UserID))
				ctx.JSON(http.StatusNotFound, response)
				return
			}

			response := helper.APIResponseError(http.StatusInternalServerError, "Create campaign failed!", err.Error())
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		req.User = user.User{ID: req.UserID}
	} else {
		req.User = userData
	}

	newCampaignData, err := handler.campaignSvc.CreateCampaign(req)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Create campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignData(newCampaignData)
	response := helper.APIResponse(http.StatusCreated, "Create campaign successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v creating campaign id %v.", userData.Name, newCampaignData.ID))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *campaignHandler) UpdateCampaign(ctx *gin.Context) {
	var reqID campaign.RequestGetCampaignByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqUpdate campaign.RequestUpdateCampaign

	err = ctx.ShouldBind(&reqUpdate)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqUpdate.User = ctx.MustGet("userData").(user.User)

	oldCampaign, err := handler.campaignSvc.GetCampaignByID(reqID)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Update campaign failed!", fmt.Sprintf("Campaign with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Update campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	updatedCampaign, err := handler.campaignSvc.UpdateCampaign(reqID, reqUpdate)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Update campaign failed!", fmt.Sprintf("Campaign with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Update campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if oldCampaign.Status != "active" && updatedCampaign.Status == "active" {
		ownerCampaignUserData, err := handler.userSvc.GetUserByID(updatedCampaign.UserID)

		if err != nil {
			response := helper.APIResponseError(http.StatusInternalServerError, "Update campaign failed!", err.Error())
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		{
			templateData := helper.EmailCampaignActive{
				Name:         ownerCampaignUserData.Name,
				Campaign:     updatedCampaign,
				GoalAmount:   helper.FormatRupiah(float64(updatedCampaign.GoalAmount)),
				CampaignLink: os.Getenv("WEB_URL") + "/donate/" + strconv.Itoa(updatedCampaign.ID),
			}
			go helper.SendMail(ownerCampaignUserData.Email, "Your Donation Campaign Now Active!", templateData, "html/campaign_active.html")
		}
	}

	formatData := campaign.FormatCampaignData(updatedCampaign)
	response := helper.APIResponse(http.StatusOK, "Update campaign successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v updating campaign id %v.", reqUpdate.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) DeleteCampaign(ctx *gin.Context) {
	var reqID campaign.RequestGetCampaignByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete campaign.RequestDeleteCampaign

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.campaignSvc.DeleteCampaign(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete campaign failed!", fmt.Sprintf("Campaign with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete campaign successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting campaign id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) GetAllCampaignImage(ctx *gin.Context) {
	campaignImages, err := handler.campaignSvc.GetAllCampaignImage()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get campaign images failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatMultipleCampaignImageData(campaignImages)
	response := helper.APIResponse(http.StatusOK, "Get campaign images successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) GetCampaignImageByID(ctx *gin.Context) {
	var req campaign.RequestGetCampaignImageByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get detail campaign image failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	campaignImageDetail, err := handler.campaignSvc.GetCampaignImageByID(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get detail campaign image failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get detail campaign image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignImageData(campaignImageDetail)
	response := helper.APIResponse(http.StatusOK, "Get detail campaign image successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) UploadImage(ctx *gin.Context) {
	var req campaign.RequestCreateCampaignImage

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Upload campaign image failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	req.User = ctx.MustGet("userData").(user.User)

	file, err := ctx.FormFile("file")
	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Upload campaign image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	slug := slug.Make(fmt.Sprintf("%d %v %s", 1, time.Now().Unix(), file.Filename[:len(file.Filename)-len(filepath.Ext(file.Filename))]))
	path := fmt.Sprintf("images/%s%v", slug, filepath.Ext(file.Filename))

	if err := ctx.SaveUploadedFile(file, path); err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Upload campaign image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	uploadedCampaignImage, err := handler.campaignSvc.SaveCampaignImage(req, path)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Upload campaign image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignImageData(uploadedCampaignImage)
	response := helper.APIResponse(http.StatusOK, "Upload campaign image successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v uploading image id %v for campaign id %v.", req.User.Name, uploadedCampaignImage.ID, uploadedCampaignImage.CampaignID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) DeleteCampaignImage(ctx *gin.Context) {
	var reqID campaign.RequestGetCampaignImageByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete campaign image failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete campaign.RequestDeleteCampaignImage

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete campaign image failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.campaignSvc.DeleteCampaignImage(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete campaign image failed!", fmt.Sprintf("Campaign image with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete campaign image failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete campaign image successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting campaign image id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) GetAllCampaignCategory(ctx *gin.Context) {
	campaignCategories, err := handler.campaignSvc.GetAllCampaignCategory()

	if err != nil {
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get campaign categories failed!", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatData := campaign.FormatMultipleCampaignCategoryData(campaignCategories)
	response := helper.APIResponse(http.StatusOK, "Get campaign categories successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) GetCampaignCategoryByID(ctx *gin.Context) {
	var req campaign.RequestGetCampaignCategoryByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get detail campaign category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	campaignCategoryDetail, err := handler.campaignSvc.GetCampaignCategoryByID(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get detail campaign category failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get detail campaign category failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignCategoryData(campaignCategoryDetail)
	response := helper.APIResponse(http.StatusOK, "Get detail campaign category successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) DeleteCampaignCategory(ctx *gin.Context) {
	var reqID campaign.RequestGetCampaignCategoryByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete campaign category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete campaign.RequestDeleteCampaignCategory

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete campaign category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.campaignSvc.DeleteCampaignCategory(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete campaign category failed!", fmt.Sprintf("Campaign category with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete campaign category failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete campaign category successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting category id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) CreateCampaignCategory(ctx *gin.Context) {
	var req campaign.RequestCreateCampaignCategory

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Create campaign category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	req.User = ctx.MustGet("userData").(user.User)

	newCampaignCategoryData, err := handler.campaignSvc.CreateCampaignCategory(req)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Create campaign category failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignCategoryData(newCampaignCategoryData)
	response := helper.APIResponse(http.StatusCreated, "Create campaign category successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v creating category id %v.", req.User.Name, newCampaignCategoryData.ID))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *campaignHandler) UpdateCampaignCategory(ctx *gin.Context) {
	var reqID campaign.RequestGetCampaignCategoryByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update campaign category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqUpdate campaign.RequestUpdateCampaignCategory

	err = ctx.ShouldBind(&reqUpdate)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update campaign category failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqUpdate.User = ctx.MustGet("userData").(user.User)

	updatedCampaignCategory, err := handler.campaignSvc.UpdateCampaignCategory(reqID, reqUpdate)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Update campaign category failed!", fmt.Sprintf("Campaign category with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Update campaign category failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignCategoryData(updatedCampaignCategory)
	response := helper.APIResponse(http.StatusOK, "Update campaign category successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v updating category id %v.", reqUpdate.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) GetAllCampaignExclusive(ctx *gin.Context) {
	exclusiveCampaigns, err := handler.campaignSvc.GetAllCampaignExclusive()

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get exclusive campaigns failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatMultipleCampaignExclusiveData(exclusiveCampaigns)
	response := helper.APIResponse(http.StatusOK, "Get exclusive campaigns successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) CreateCampaignExclusive(ctx *gin.Context) {
	var req campaign.RequestCreateCampaignExclusive

	err := ctx.ShouldBind(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Create exclusive campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	req.User = ctx.MustGet("userData").(user.User)

	newCampaignExclusiveData, err := handler.campaignSvc.CreateCampaignExclusive(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Create exclusive campaign failed!", fmt.Sprintf("Campaign with ID %d not found!", req.CampaignID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Create exclusive campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignExclusiveData(newCampaignExclusiveData)
	response := helper.APIResponse(http.StatusCreated, "Create exclusive campaign successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v set campaign id %v to exclusive.", req.User.Name, req.CampaignID))

	ctx.JSON(http.StatusCreated, response)
}

func (handler *campaignHandler) GetCampaignExclusiveByID(ctx *gin.Context) {
	var req campaign.RequestGetCampaignExclusiveByID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get detail exclusive campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	exclusiveCampaignDetail, err := handler.campaignSvc.GetCampaignExclusiveByID(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get detail exclusive campaign failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get detail exclusive campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignExclusiveData(exclusiveCampaignDetail)
	response := helper.APIResponse(http.StatusOK, "Get detail exclusive campaign successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) GetCampaignExclusiveByCampaignID(ctx *gin.Context) {
	var req campaign.RequestGetCampaignExclusiveByCampaignID

	err := ctx.ShouldBindUri(&req)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get detail exclusive campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	exclusiveCampaignDetail, err := handler.campaignSvc.GetCampaignExclusiveByCampaignID(req)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Get detail exclusive campaign failed!", "Data not found!")
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Get detail exclusive campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	formatData := campaign.FormatCampaignExclusiveData(exclusiveCampaignDetail)
	response := helper.APIResponse(http.StatusOK, "Get detail exclusive campaign successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) GetCampaignExclusiveByWinnerUserID(ctx *gin.Context) {
	var req campaign.RequestGetCampaignExclusiveByWinnerUserID

	req.ID = ctx.MustGet("userData").(user.User).ID

	exclusiveCampaigns, err := handler.campaignSvc.GetCampaignExclusiveByWinnerUserID(req)

	if err != nil {
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Get exclusive campaigns failed!", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatData := campaign.FormatMultipleCampaignExclusiveData(exclusiveCampaigns)
	response := helper.APIResponse(http.StatusOK, "Get exclusive campaigns successfully!", formatData)

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) UpdateCampaignExclusive(ctx *gin.Context) {
	var reqID campaign.RequestGetCampaignExclusiveByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update exclusive campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqUpdate campaign.RequestUpdateCampaignExclusive

	err = ctx.ShouldBind(&reqUpdate)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Update exclusive campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqUpdate.User = ctx.MustGet("userData").(user.User)

	updatedCampaignExclusive, err := handler.campaignSvc.UpdateCampaignExclusive(reqID, reqUpdate)

	if err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Update exclusive campaign failed!", fmt.Sprintf("Exclusive campaign with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Update exclusive campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if updatedCampaignExclusive.WinnerUserID != 0 {
		winnerUserData, err := handler.userSvc.GetUserByID(updatedCampaignExclusive.WinnerUserID)

		if err != nil {
			response := helper.APIResponseError(http.StatusInternalServerError, "Update exclusive campaign failed!", err.Error())
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}

		{
			status := "Pending"
			reward := updatedCampaignExclusive.Reward

			if updatedCampaignExclusive.IsPaidOff == 1 {
				status = "Paid Off"
			}

			if updatedCampaignExclusive.IsRewardMoney == 1 {
				convertedReward, _ := strconv.Atoi(updatedCampaignExclusive.Reward)
				reward = helper.FormatRupiah(float64(convertedReward))
			}

			templateData := helper.EmailRewardUpdate{
				CampaignLink: os.Getenv("WEB_URL") + "/donate/" + strconv.Itoa(updatedCampaignExclusive.CampaignID),
				Name:         winnerUserData.Name,
				Reward:       reward,
				Status:       status,
			}
			go helper.SendMail(winnerUserData.Email, "Information Update For Your Reward", templateData, "html/reward_update.html")
		}
	}

	formatData := campaign.FormatCampaignExclusiveData(updatedCampaignExclusive)
	response := helper.APIResponse(http.StatusOK, "Update exclusive campaign successfully!", formatData)

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v updating campaign exclusive id %v.", reqUpdate.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) DeleteCampaignExclusive(ctx *gin.Context) {
	var reqID campaign.RequestGetCampaignExclusiveByID

	err := ctx.ShouldBindUri(&reqID)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete exclusive campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	var reqDelete campaign.RequestDeleteCampaignExclusive

	err = ctx.ShouldBind(&reqDelete)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponseError(http.StatusUnprocessableEntity, "Delete exclusive campaign failed!", errors[0])
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	reqDelete.User = ctx.MustGet("userData").(user.User)

	if _, err = handler.campaignSvc.DeleteCampaignExclusive(reqID, reqDelete); err != nil {
		if helper.IsErrNoRows(err.Error()) {
			response := helper.APIResponseError(http.StatusNotFound, "Delete exclusive campaign failed!", fmt.Sprintf("Exclusive Campaign with ID %d not found!", reqID.ID))
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		response := helper.APIResponseError(http.StatusInternalServerError, "Delete exclusive campaign failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BasicAPIResponse(http.StatusOK, "Delete exclusive campaign successfully!")

	handler.logsSvc.CreateActivityLog(ctx, fmt.Sprintf("%v deleting campaign exclusive id %v.", reqDelete.User.Name, reqID.ID))

	ctx.JSON(http.StatusOK, response)
}

func (handler *campaignHandler) AdminDataTablesCampaigns(ctx *gin.Context) {
	dataTablesCampaigns, err := handler.campaignSvc.AdminDataTablesCampaigns(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables campaigns failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesCampaigns)
}

func (handler *campaignHandler) AdminDataTablesCategories(ctx *gin.Context) {
	dataTablesCategories, err := handler.campaignSvc.AdminDataTablesCategories(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables categories failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesCategories)
}

func (handler *campaignHandler) UserDataTablesCampaigns(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(user.User)
	dataTablesCampaigns, err := handler.campaignSvc.UserDataTablesCampaigns(ctx, userData)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables campaigns failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesCampaigns)
}

func (handler *campaignHandler) AdminDataTablesWinnersExclusiveCampaigns(ctx *gin.Context) {
	dataTablesWinnerExclusiveCampaigns, err := handler.campaignSvc.AdminDataTablesWinnersExclusiveCampaigns(ctx)

	if err != nil {
		response := helper.APIResponseError(http.StatusInternalServerError, "Get datatables winners exclusive campaigns failed!", err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusOK, dataTablesWinnerExclusiveCampaigns)
}
