package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/WeAreAmazingTeam/tcd-backend/auth"
	"github.com/WeAreAmazingTeam/tcd-backend/campaign"
	"github.com/WeAreAmazingTeam/tcd-backend/chart"
	"github.com/WeAreAmazingTeam/tcd-backend/company"
	theCloudConfig "github.com/WeAreAmazingTeam/tcd-backend/config"
	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/WeAreAmazingTeam/tcd-backend/handler"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/logs"
	"github.com/WeAreAmazingTeam/tcd-backend/middleware"
	"github.com/WeAreAmazingTeam/tcd-backend/payment"
	"github.com/WeAreAmazingTeam/tcd-backend/transaction"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	isProduction := flag.Bool("production", true, "production mode?")
	flag.Parse()

	_, b, _, _ := runtime.Caller(0)
	projectRootPath := filepath.Join(filepath.Dir(b), "")
	envLocation := projectRootPath + "/.env"

	if *isProduction {
		envLocation = "/www/wwwroot/golang/.env"
	}

	if err := godotenv.Load(envLocation); err != nil {
		log.Fatal("error while loading or open .env file, err: ", err.Error())
	}

	// initial constants
	constant.InitDBConstant()
	constant.InitAuthConstant()
	constant.InitRedisConstant()

	// initial database
	db := theCloudConfig.InitDB(*isProduction)

	// initial scheduler
	theCloudConfig.InitScheduler(db)

	// repositories
	userRepository := user.NewRepository(db)
	chartRepository := chart.NewRepository(db)
	companyRepository := company.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)
	logsRepository := logs.NewRepository(db)

	// services
	userSvc := user.NewService(userRepository)
	authSvc := auth.NewService()
	chartSvc := chart.NewService(chartRepository)
	paymentSvc := payment.NewService()
	campaignSvc := campaign.NewService(campaignRepository, userRepository, companyRepository)
	companySvc := company.NewService(companyRepository)
	transactionSvc := transaction.NewService(transactionRepository, campaignRepository, userRepository, companyRepository, campaignSvc, paymentSvc)
	logsSvc := logs.NewService(logsRepository)

	// handlers
	userHandler := handler.NewUserHandler(userSvc, authSvc, logsSvc, companySvc)
	chartHandler := handler.NewChartHandler(chartSvc)
	campaignHandler := handler.NewCampaignHandler(campaignSvc, userSvc, logsSvc)
	companyHandler := handler.NewCompanyHandler(companySvc, logsSvc)
	transactionHandler := handler.NewTransactionHandler(transactionSvc, campaignSvc, paymentSvc, userSvc, logsSvc)
	logsHandler := handler.NewLogsHandler(logsSvc)
	webAndCMSHandler := handler.NewWebAndCMSHandler(transactionSvc, campaignSvc, paymentSvc, userSvc, logsSvc)

	// for activate release mode
	if *isProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	// gin app configuration
	app := gin.Default()
	app.SetTrustedProxies(nil)
	app.Static("/images", "./images")
	app.Use(gzip.Gzip(gzip.DefaultCompression))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "POST"},
		AllowHeaders:     []string{"Host", "Origin", "Content-Length", "Content-Type", "Authorization", "User-Agent", "X-Forwarded-For", "Accept-Encoding", "Connection"},
		ExposeHeaders:    []string{"Content-Length", "Content-Encoding"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// middleware
	mAuth := middleware.Auth(authSvc, userSvc)
	mAdminAuth := middleware.AdminAuth(authSvc, userSvc)

	// routing
	api := app.Group("/api/v1")
	{
		// >>>>>>>>>>>>>>> begin strict endpoint <<<<<<<<<<<<<<<

		// account settings
		api.GET("/users/data", mAuth, userHandler.GetUserData)
		api.PUT("/users/data/change", mAuth, userHandler.ChangeUserData)
		api.POST("/users/withdraw", mAuth, userHandler.CreateWithdrawalRequest)

		// users (for admin only)
		api.GET("/users", mAdminAuth, userHandler.GetAllUser)
		api.GET("/users/:id", mAdminAuth, userHandler.GetUserByID)
		api.PUT("/users/withdrawal/:id", mAdminAuth, userHandler.UpdateUserWithdrawalRequest)
		api.PUT("/users/:id", mAdminAuth, userHandler.UpdateUser)
		api.POST("/users", mAdminAuth, userHandler.CreateUser)
		api.DELETE("/users/withdrawal/:id", mAdminAuth, userHandler.DeleteUserWithdrawalRequest)
		api.DELETE("/users/:id", mAdminAuth, userHandler.DeleteUser)

		// campaigns
		api.PUT("/campaigns/:id", mAuth, campaignHandler.UpdateCampaign)
		api.POST("/campaigns", mAuth, campaignHandler.CreateCampaign)
		api.DELETE("/campaigns/:id", mAuth, campaignHandler.DeleteCampaign)

		// campaigns -> images
		api.POST("/campaigns/images", mAuth, campaignHandler.UploadImage)
		api.DELETE("/campaigns/images/:id", mAuth, campaignHandler.DeleteCampaignImage)

		// campaigns -> categories (for admin only)
		api.PUT("/campaigns/categories/:id", mAdminAuth, campaignHandler.UpdateCampaignCategory)
		api.POST("/campaigns/categories", mAdminAuth, campaignHandler.CreateCampaignCategory)
		api.DELETE("/campaigns/categories/:id", mAdminAuth, campaignHandler.DeleteCampaignCategory)

		// for get exclusive campaign by user id
		api.GET("/campaigns/exclusive/user", mAuth, campaignHandler.GetCampaignExclusiveByWinnerUserID)

		// campaigns exclusive (for admin only)
		api.GET("/campaigns/exclusive", mAdminAuth, campaignHandler.GetAllCampaignExclusive)
		api.GET("/campaigns/exclusive/:id", mAdminAuth, campaignHandler.GetCampaignExclusiveByID)
		api.PUT("/campaigns/exclusive/:id", mAdminAuth, campaignHandler.UpdateCampaignExclusive)
		api.POST("/campaigns/exclusive", mAdminAuth, campaignHandler.CreateCampaignExclusive)
		api.DELETE("/campaigns/exclusive/:id", mAdminAuth, campaignHandler.DeleteCampaignExclusive)

		// transactions
		api.POST("/transactions", mAuth, transactionHandler.CreateTransaction)
		api.POST("/transactions/emoney", mAuth, transactionHandler.CreateTransactionWithEMoney)
		api.DELETE("/transactions/:id", mAuth, transactionHandler.DeleteTransaction)

		// company -> cash flow
		api.POST("/company/cashflow", mAdminAuth, companyHandler.CreateCompanyCashFlow)
		api.DELETE("/company/cashflow/:id", mAdminAuth, companyHandler.DeleteCompanyCashFlow)

		// admin datatables
		api.GET("admin/datatables/users", mAdminAuth, userHandler.AdminDataTablesUsers)
		api.GET("admin/datatables/categories", mAdminAuth, campaignHandler.AdminDataTablesCategories)
		api.GET("admin/datatables/campaigns", mAdminAuth, campaignHandler.AdminDataTablesCampaigns)
		api.GET("admin/datatables/transactions", mAdminAuth, transactionHandler.AdminDataTablesTransactions)
		api.GET("admin/datatables/logs/activity", mAdminAuth, logsHandler.AdminDataTablesActivityLogs)
		api.GET("admin/datatables/campaigns/exclusive", mAdminAuth, campaignHandler.AdminDataTablesWinnersExclusiveCampaigns)
		api.GET("admin/datatables/withdrawal", mAdminAuth, userHandler.AdminDatatablesWithdrawalRequest)
		api.GET("admin/datatables/company/cashflow", mAdminAuth, companyHandler.AdminDataTablesCompanyCashFlow)

		// datatables for user
		api.GET("datatables/campaigns", mAuth, campaignHandler.UserDataTablesCampaigns)
		api.GET("datatables/transactions", mAuth, transactionHandler.UserDataTablesTransactions)
		api.GET("datatables/flow/emoney", mAuth, userHandler.UserDataTablesEMoneyFlow)
		api.GET("datatables/withdrawal", mAuth, userHandler.UserDatatablesWithdrawalRequest)

		// logs
		api.POST("logs/activity/auth", mAuth, logsHandler.AddLogsActivityAuth)

		// dashboard statistics
		api.GET("admin/dashboard/statistics", mAdminAuth, webAndCMSHandler.GetStatisticsForAdminDashboard)

		// get chart
		api.GET("admin/dashboard/chart", mAdminAuth, chartHandler.GetChart)

		// >>>>>>>>>>>>>>> end strict endpoint <<<<<<<<<<<<<<<

		// >>>>>>>>>>>>>>> begin non-strict endpoint <<<<<<<<<<<<<<<

		// authentication
		api.POST("/users/register", userHandler.Register)
		api.POST("/users/login", userHandler.Login)

		// forgot password
		api.GET("/users/forgot-password/:token", userHandler.ProcessForgotPasswordToken)
		api.POST("/users/forgot-password", userHandler.CreateForgotPasswordToken)

		// users
		api.GET("/users/name/:id", userHandler.GetNameByID)

		// campaigns
		api.GET("/campaigns", campaignHandler.GetAllCampaign)
		api.GET("/campaigns/:id", campaignHandler.GetCampaignByID)

		// campaigns -> images
		api.GET("/campaigns/images", campaignHandler.GetAllCampaignImage)
		api.GET("/campaigns/images/:id", campaignHandler.GetCampaignImageByID)

		// campaigns -> categories
		api.GET("/campaigns/categories", campaignHandler.GetAllCampaignCategory)
		api.GET("/campaigns/categories/:id", campaignHandler.GetCampaignCategoryByID)

		// campaigns exclusive by campaign id
		api.GET("/campaigns/exclusive/campaign/:id", campaignHandler.GetCampaignExclusiveByCampaignID)

		// transactions
		api.GET("/transactions", transactionHandler.GetAllTransaction)
		api.GET("/transactions/:id", transactionHandler.GetTransactionByID)
		api.GET("/transactions/users/:id", transactionHandler.GetTransactionByUserID)
		api.GET("/transactions/campaigns/:id", transactionHandler.GetTransactionByCampaignID)
		api.POST("/transactions/test/midtrans", transactionHandler.TestMidtrans)
		api.POST("/transactions/webhooks", transactionHandler.TransactionWebhooks)
		api.POST("/transactions/anonymous", transactionHandler.CreateAnonymousTransaction)

		// logs
		api.POST("logs/activity", logsHandler.AddLogsActivity)

		// web
		api.GET("web/home/statistics", webAndCMSHandler.GetStatisticsForHomePage)

		// >>>>>>>>>>>>>>> end non-strict endpoint <<<<<<<<<<<<<<<
	}

	// handle invalid method
	app.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, helper.BasicAPIResponseError(http.StatusMethodNotAllowed, "Request invalid, invalid method!"))
	})

	// handle invalid path or invalid endpoint
	app.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, helper.BasicAPIResponseError(http.StatusNotFound, "Request invalid, path not found!"))
	})

	// run http server
	app.Run(os.Getenv("APP_RUN_ON"))
}
