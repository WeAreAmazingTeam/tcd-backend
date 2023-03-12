package campaign

import (
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllCampaign(ctx *gin.Context) ([]Campaign, error)
	GetCampaignByID(id int) (Campaign, error)
	SaveCampaign(Campaign) (Campaign, error)
	UpdateCampaign(Campaign) (Campaign, error)
	UpdateCampaignFromPayment(campaignID int, transactionAmount int64) error
	DeleteCampaign(Campaign) (bool, error)

	GetAllCampaignImage() ([]CampaignImage, error)
	GetCampaignImageByID(id int) (CampaignImage, error)
	CreateCampaignImage(CampaignImage) (CampaignImage, error)
	UpdateAllImagesAsNonPrimary(campaignID int) (bool, error)
	DeleteCampaignImage(CampaignImage) (bool, error)

	GetAllCampaignCategory() ([]CampaignCategory, error)
	GetCampaignCategoryByID(id int) (CampaignCategory, error)
	SaveCampaignCategory(CampaignCategory) (CampaignCategory, error)
	UpdateCampaignCategory(CampaignCategory) (CampaignCategory, error)
	DeleteCampaignCategory(CampaignCategory) (bool, error)

	GetAllCampaignExclusive() ([]ExclusiveCampaign, error)
	GetCampaignExclusiveByID(id int) (ExclusiveCampaign, error)
	GetCampaignExclusiveByCampaignID(id int) (ExclusiveCampaign, error)
	GetCampaignExclusiveByWinnerUserID(id int) ([]ExclusiveCampaign, error)
	SaveCampaignExclusive(ExclusiveCampaign) (ExclusiveCampaign, error)
	UpdateCampaignExclusive(ExclusiveCampaign) (ExclusiveCampaign, error)
	GetWinnerCampaignExclusive(ExclusiveCampaign) (int, error)
	DeleteCampaignExclusive(ExclusiveCampaign) (bool, error)

	AdminDataTablesCampaigns(ctx *gin.Context) (helper.DataTables, error)
	AdminDataTablesCategories(ctx *gin.Context) (helper.DataTables, error)
	AdminDataTablesWinnersExclusiveCampaigns(*gin.Context) (helper.DataTables, error)

	UserDataTablesCampaigns(*gin.Context, user.User) (helper.DataTables, error)

	GetTotalDonation() (int, error)
	GetDonationCompleted() (int, error)
}

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{DB: db}
}
