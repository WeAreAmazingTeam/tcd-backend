package helper

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"gopkg.in/gomail.v2"
)

type EmailWelcome struct {
	Name string
}

type EmailForgotPassword struct {
	Name string
	URL  string
}

type EmailCampaignFinished struct {
	Campaign       any
	Name           string
	GoalAmount     string
	CollectedFunds string
	AdminFee       string
	FinalAmount    string
}

type EmailEarningRewardFromExclusiveCampaign struct {
	CampaignLink string
	Name         string
	Reward       string
	Status       string
}

type EmailTransactionSuccess struct {
	CampaignLink string
	Name         string
	Amount       string
}

type EmailWithdrawalRequest struct {
	Name   string
	Amount string
}

type EmailWithdrawalApproved struct {
	Name   string
	Amount string
}

type EmailWithdrawalRejected struct {
	Name   string
	Amount string
}

type EmailRewardUpdate struct {
	CampaignLink string
	Name         string
	Reward       string
	Status       string
}

type EmailCampaignActive struct {
	Name         string
	Campaign     any
	GoalAmount   string
	CampaignLink string
}

func ParseTemplate(templateFileName string, data any) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SendMail(to string, subject string, data any, templateFile string) {
	notYetSuccess := true
	html, err := ParseTemplate(templateFile, data)

	if err != nil {
		fmt.Printf("[MAIL] email failed to send to %v while parse HTML template, template: %v, [err: %v]\n", to, templateFile, err.Error())
		return
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("MAIL_ADDR_HEADER"))
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", html)

	p, _ := strconv.Atoi(os.Getenv("MAIL_SMTP_PORT"))
	dial := gomail.NewDialer(os.Getenv("MAIL_SMTP_HOST"), p, os.Getenv("MAIL_ADDR"), os.Getenv("MAIL_PASS"))

	if err := dial.DialAndSend(mailer); err == nil {
		fmt.Printf("[MAIL] email sent successfully to %v, subject: %v, template: %v.\n", to, subject, templateFile)
		notYetSuccess = false
	} else {
		fmt.Printf("[MAIL] email failed to send to %v while dial and send, [err: %v]\n", to, err.Error())
	}

	if notYetSuccess {
		mailer := gomail.NewMessage()
		mailer.SetHeader("From", os.Getenv("BACKUP_MAIL_ADDR_HEADER"))
		mailer.SetHeader("To", to)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/html", html)

		p, _ := strconv.Atoi(os.Getenv("BACKUP_MAIL_SMTP_PORT"))
		dial := gomail.NewDialer(os.Getenv("BACKUP_MAIL_SMTP_HOST"), p, os.Getenv("BACKUP_MAIL_ADDR"), os.Getenv("BACKUP_MAIL_PASS"))

		if err := dial.DialAndSend(mailer); err == nil {
			fmt.Printf("[MAIL FIRST BACKUP] email sent successfully to %v, subject: %v, template: %v.\n", to, subject, templateFile)
			notYetSuccess = false
		} else {
			fmt.Printf("[MAIL FIRST BACKUP] email failed to send to %v while dial and send, [err: %v]\n", to, err.Error())
		}
	}

	if notYetSuccess {
		mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_PRIKEY"))
		msg := mg.NewMessage(os.Getenv("MAILGUN_SENDER"), subject, "", to)

		msg.SetHtml(html)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		resp, id, err := mg.Send(ctx, msg)

		if err != nil {
			fmt.Printf("[MAIL LAST BACKUP] email failed to send to %v, [err: %v]\n", to, err.Error())
		} else {
			fmt.Printf("[MAIL LAST BACKUP] email sent successfully to %v, subject: %v, template: %v. [ID: %s Resp: %s]\n", to, subject, templateFile, id, resp)
		}
	}
}
