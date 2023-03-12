package config

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/WeAreAmazingTeam/tcd-backend/campaign"
	"github.com/WeAreAmazingTeam/tcd-backend/company"
	"github.com/WeAreAmazingTeam/tcd-backend/helper"
	"github.com/WeAreAmazingTeam/tcd-backend/logs"
	"github.com/WeAreAmazingTeam/tcd-backend/user"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func InitScheduler(db *gorm.DB) {
	jakartaTime, err := time.LoadLocation("Asia/Jakarta")

	if err != nil {
		log.Fatal("error while load time location, err: ", err.Error())
	}

	scheduler := cron.New(cron.WithLocation(jakartaTime))

	defer scheduler.Stop()

	// for testing: */1 * * * *
	// for prod: 10 0 * * *
	scheduler.AddFunc("10 0 * * *", func() {
		affected := 0
		rows, err := db.Raw(helper.ConvertToInLineQuery(campaign.QueryGetAll+"AND status = 'active' AND finished_at <= ?"), fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05"))).Rows()
		activityLog := logs.ActivityLog{}
		activityLog.IpAddress = "-"
		activityLog.UserAgent = "-"

		if err != nil {
			activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 1)] %v", err.Error())

			log.Println(activityLog.Content)

			if err := db.Create(&activityLog).Error; err != nil {
				log.Fatal(err.Error())
			}
		}

		defer rows.Close()

		for rows.Next() {
			tmp := campaign.Campaign{}
			err := rows.Scan(
				&tmp.ID,
				&tmp.UserID,
				&tmp.CategoryID,
				&tmp.Title,
				&tmp.Slug,
				&tmp.ShortDescription,
				&tmp.Description,
				&tmp.GoalAmount,
				&tmp.CurrentAmount,
				&tmp.IsExclusive,
				&tmp.DonorCount,
				&tmp.Status,
				&tmp.FinishedAt,
				&tmp.CreatedAt,
				&tmp.CreatedBy,
				&tmp.UpdatedAt,
				&tmp.UpdatedBy,
				&tmp.DeletedAt,
				&tmp.DeletedBy,
			)

			if err != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 2)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			currentAmount := float64(tmp.CurrentAmount)
			forDeducted := currentAmount - math.Round(float64(currentAmount-(currentAmount*(float64(6)/float64(100)))))
			deductedAmount := currentAmount - forDeducted

			result := db.Model(&user.User{}).Where("id = ?", tmp.UserID).Update("e_money", gorm.Expr("e_money + ?", deductedAmount))

			if result.Error != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 3-1)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			userEMoneyFlow := user.UserEMoneyFlow{
				UserID: tmp.UserID,
				Status: "in",
				Amount: int64(currentAmount),
				Note:   fmt.Sprintf("Funds from the donation campaign: %v.", tmp.Title),
			}

			if err := db.Create(&userEMoneyFlow).Error; err != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 3-2)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			userEMoneyFlow = user.UserEMoneyFlow{
				UserID: tmp.UserID,
				Status: "out",
				Amount: int64(deductedAmount),
				Note:   fmt.Sprintf("Admin fee for the donation campaign: %v.", tmp.Title),
			}

			if err := db.Create(&userEMoneyFlow).Error; err != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 3-3)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			companyCashFlow := company.CompanyCashFlow{
				Status: "out",
				Amount: int64(currentAmount),
				Note:   fmt.Sprintf("Disburse funds for exclusive campaign: %v.", tmp.Title),
			}

			if err := db.Create(&companyCashFlow).Error; err != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 3-4)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			companyCashFlow = company.CompanyCashFlow{
				Status: "in",
				Amount: int64(deductedAmount),
				Note:   fmt.Sprintf("Admin fee from donation campaign: %v.", tmp.Title),
			}

			if err := db.Create(&companyCashFlow).Error; err != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 3-5)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			result = db.Model(&campaign.Campaign{}).Where("id = ?", tmp.ID).Update("status", "finished")

			if result.Error != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 4-1)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			var userData user.User

			if err := db.Where("id = ?", tmp.UserID).Find(&userData).Error; err != nil {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 4-2)] %v", err.Error())

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			}

			if userData.ID == 0 {
				activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 4-3)] %v", "sql: no rows in result set")

				log.Println(activityLog.Content)

				if err := db.Create(&activityLog).Error; err != nil {
					log.Fatal(err.Error())
				}
			} else {
				templateData := helper.EmailCampaignFinished{
					Campaign:       tmp,
					Name:           userData.Name,
					GoalAmount:     helper.FormatRupiah(float64(tmp.GoalAmount)),
					CollectedFunds: helper.FormatRupiah(float64(tmp.CurrentAmount)),
					AdminFee:       helper.FormatRupiah(forDeducted),
					FinalAmount:    helper.FormatRupiah(deductedAmount),
				}
				go helper.SendMail(userData.Email, fmt.Sprintf("Your Campaign (%v) Has Finished", tmp.Title), templateData, "html/campaign_finished.html")
			}

			if tmp.IsExclusive == 1 {
				exclusiveCampaign := campaign.ExclusiveCampaign{}
				rowExclusiveCampaign := db.Raw(helper.ConvertToInLineQuery(campaign.QueryGetCampaignExclusiveByCampaignID), tmp.ID).Row()

				err = rowExclusiveCampaign.Scan(
					&exclusiveCampaign.ID,
					&exclusiveCampaign.CampaignID,
					&exclusiveCampaign.WinnerUserID,
					&exclusiveCampaign.IsRewardMoney,
					&exclusiveCampaign.Reward,
					&exclusiveCampaign.IsPaidOff,
					&exclusiveCampaign.CreatedAt,
					&exclusiveCampaign.CreatedBy,
					&exclusiveCampaign.UpdatedAt,
					&exclusiveCampaign.UpdatedBy,
					&exclusiveCampaign.DeletedAt,
					&exclusiveCampaign.DeletedBy,
				)

				if err != nil {
					activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 5)] %v", err.Error())

					log.Println(activityLog.Content)

					if err := db.Create(&activityLog).Error; err != nil {
						log.Fatal(err.Error())
					}
				}

				if exclusiveCampaign.IsPaidOff == 0 {
					winnerUserID := 0
					err := db.Raw(helper.ConvertToInLineQuery(campaign.QueryGetOneRandomUserIDForWinnerExclusiveCampaign), exclusiveCampaign.CampaignID).Row().Scan(&winnerUserID)

					if err != nil {
						activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 6)] %v", err.Error())

						log.Println(activityLog.Content)

						if err := db.Create(&activityLog).Error; err != nil {
							log.Fatal(err.Error())
						}
					}

					if winnerUserID == 0 {
						activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 7)] %v", "no user can be the winner")

						log.Println(activityLog.Content)

						if err := db.Create(&activityLog).Error; err != nil {
							log.Fatal(err.Error())
						}
					} else {
						exclusiveCampaign.WinnerUserID = winnerUserID

						if exclusiveCampaign.IsRewardMoney == 1 {
							exclusiveCampaign.IsPaidOff = 1
						}

						if err := db.Save(&exclusiveCampaign).Error; err != nil {
							activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 8)] %v", err.Error())

							log.Println(activityLog.Content)

							if err := db.Create(&activityLog).Error; err != nil {
								log.Fatal(err.Error())
							}
						}

						if exclusiveCampaign.IsRewardMoney == 1 {
							var userData user.User

							if err := db.Where("id = ?", winnerUserID).Find(&userData).Error; err != nil {
								activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 9)] %v", err.Error())

								log.Println(activityLog.Content)

								if err := db.Create(&activityLog).Error; err != nil {
									log.Fatal(err.Error())
								}
							}

							moneyReward, err := strconv.Atoi(exclusiveCampaign.Reward)

							if err != nil {
								activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 10)] %v", err.Error())

								log.Println(activityLog.Content)

								if err := db.Create(&activityLog).Error; err != nil {
									log.Fatal(err.Error())
								}
							}

							userData.EMoney = userData.EMoney + float64(moneyReward)

							if err := db.Save(&userData).Error; err != nil {
								activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 11-1)] %v", err.Error())

								log.Println(activityLog.Content)

								if err := db.Create(&activityLog).Error; err != nil {
									log.Fatal(err.Error())
								}
							}

							userEMoneyFlow := user.UserEMoneyFlow{
								UserID: userData.ID,
								Status: "in",
								Amount: int64(moneyReward),
								Note:   fmt.Sprintf("Reward from exclusive campaign id %v.", exclusiveCampaign.CampaignID),
							}

							if err := db.Create(&userEMoneyFlow).Error; err != nil {
								activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 11-2)] %v", err.Error())

								log.Println(activityLog.Content)

								if err := db.Create(&activityLog).Error; err != nil {
									log.Fatal(err.Error())
								}
							}

							companyCashFlow := company.CompanyCashFlow{
								Status: "out",
								Amount: int64(moneyReward),
								Note:   fmt.Sprintf("Reward for exclusive campaign id %v.", exclusiveCampaign.CampaignID),
							}

							if err := db.Create(&companyCashFlow).Error; err != nil {
								activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 11-3)] %v", err.Error())

								log.Println(activityLog.Content)

								if err := db.Create(&activityLog).Error; err != nil {
									log.Fatal(err.Error())
								}
							}
						}

						var winnerUserData user.User

						if err := db.Where("id = ?", winnerUserID).Find(&winnerUserData).Error; err != nil {
							activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 12-1)] %v", err.Error())

							log.Println(activityLog.Content)

							if err := db.Create(&activityLog).Error; err != nil {
								log.Fatal(err.Error())
							}
						}

						if winnerUserData.ID == 0 {
							activityLog.Content = fmt.Sprintf("[CRON IMPORTANT INFO (STEP 12-2)] %v", "sql: no rows in result set")

							log.Println(activityLog.Content)

							if err := db.Create(&activityLog).Error; err != nil {
								log.Fatal(err.Error())
							}
						} else {
							status := "Pending"

							if exclusiveCampaign.IsPaidOff == 1 {
								status = "Paid Off"
							}

							templateData := helper.EmailEarningRewardFromExclusiveCampaign{
								CampaignLink: os.Getenv("WEB_URL") + "/donate/" + strconv.Itoa(exclusiveCampaign.CampaignID),
								Name:         winnerUserData.Name,
								Reward:       exclusiveCampaign.Reward,
								Status:       status,
							}
							go helper.SendMail(winnerUserData.Email, "Congratulations, You Get Rewards From Exclusive Campaign!", templateData, "html/earn_reward.html")
						}
					}
				}
			}

			affected = affected + int(result.RowsAffected)
		}

		activityLog.Content = fmt.Sprintf("System running CRON for check and update finished campaign. (affected: %v)", affected)

		log.Println(activityLog.Content)

		if err := db.Create(&activityLog).Error; err != nil {
			log.Fatal(err.Error())
		}
	})

	go scheduler.Start()
}
