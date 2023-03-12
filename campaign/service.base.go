package campaign

import (
	"github.com/WeAreAmazingTeam/tcd-backend/company"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
)

type Service interface {
	GetAllCampaign(ctx *gin.Context) ([]Campaign, error)
	GetCampaignByID(RequestGetCampaignByID) (Campaign, error)
	CreateCampaign(RequestCreateCampaign) (Campaign, error)
	UpdateCampaign(RequestGetCampaignByID, RequestUpdateCampaign) (Campaign, error)
	DeleteCampaign(RequestGetCampaignByID, RequestDeleteCampaign) (bool, error)

	GetAllCampaignImage() ([]CampaignImage, error)
	GetCampaignImageByID(RequestGetCampaignImageByID) (CampaignImage, error)
	SaveCampaignImage(RequestCreateCampaignImage, string) (CampaignImage, error)
	DeleteCampaignImage(RequestGetCampaignImageByID, RequestDeleteCampaignImage) (bool, error)

	GetAllCampaignCategory() ([]CampaignCategory, error)
	GetCampaignCategoryByID(RequestGetCampaignCategoryByID) (CampaignCategory, error)
	CreateCampaignCategory(RequestCreateCampaignCategory) (CampaignCategory, error)
	UpdateCampaignCategory(RequestGetCampaignCategoryByID, RequestUpdateCampaignCategory) (CampaignCategory, error)
	DeleteCampaignCategory(RequestGetCampaignCategoryByID, RequestDeleteCampaignCategory) (bool, error)

	GetAllCampaignExclusive() ([]ExclusiveCampaign, error)
	GetCampaignExclusiveByID(RequestGetCampaignExclusiveByID) (ExclusiveCampaign, error)
	GetCampaignExclusiveByCampaignID(RequestGetCampaignExclusiveByCampaignID) (ExclusiveCampaign, error)
	GetCampaignExclusiveByWinnerUserID(RequestGetCampaignExclusiveByWinnerUserID) ([]ExclusiveCampaign, error)
	CreateCampaignExclusive(RequestCreateCampaignExclusive) (ExclusiveCampaign, error)
	UpdateCampaignExclusive(RequestGetCampaignExclusiveByID, RequestUpdateCampaignExclusive) (ExclusiveCampaign, error)
	CheckAndSetWinnerCampaignExclusive(RequestGetCampaignExclusiveByCampaignID) (ExclusiveCampaign, error)
	DeleteCampaignExclusive(RequestGetCampaignExclusiveByID, RequestDeleteCampaignExclusive) (bool, error)

	AdminDataTablesCampaigns(*gin.Context) (helper.DataTables, error)
	AdminDataTablesCategories(*gin.Context) (helper.DataTables, error)
	AdminDataTablesWinnersExclusiveCampaigns(*gin.Context) (helper.DataTables, error)

	UserDataTablesCampaigns(*gin.Context, user.User) (helper.DataTables, error)

	GetTotalDonation() (int, error)
	GetDonationCompleted() (int, error)
}

type service struct {
	repo        Repository
	userRepo    user.Repository
	companyRepo company.Repository
}

func NewService(
	repository Repository,
	userRepository user.Repository,
	companyRepository company.Repository,
) *service {
	return &service{
		repo:        repository,
		userRepo:    userRepository,
		companyRepo: companyRepository,
	}
}
