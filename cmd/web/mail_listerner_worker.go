package main

import (
	"fmt"
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	mail "github.com/xhit/go-simple-mail/v2"
	"os"
	"strings"
	"time"
)

func startWorkerListenMail() {
	go listenForMailForever()
}

func listenForMailForever() {
	for {
		sendEmailWhenRetrievedFromChannel()
	}
}

func sendEmailWhenRetrievedFromChannel() {

	message := <-AppConfig.MailChan
	sendMail(message)

}

func sendMail(mailData model.MailData) {

	server := setupMailServer()
	client := setupMailClient(server)

	email := setupEmail(mailData)
	err := email.Send(client)
	if err != nil {
		AppConfig.ErrorLog.Print(err)
	} else {
		AppConfig.InfoLog.Print("Email Sent")
	}

}

func setupEmail(mailData model.MailData) *mail.Email {

	email := mail.NewMSG()
	email.SetFrom(mailData.From).AddTo(mailData.To).SetSubject(mailData.Subject)
	email.SetBody(mail.TextHTML, setupEmailBody(mailData.TemplateName, mailData.Content, mailData.From))

	return email

}

func setupEmailBody(templateName, content, from string) string {
	if templateName == "" {
		return content
	}

	filePath := fmt.Sprintf("./email-templates/%s", templateName)
	data, err := os.ReadFile(filePath)

	if err != nil {
		AppConfig.ErrorLog.Print(err)
	}

	mailTemplate := string(data)
	replacedContent := strings.Replace(mailTemplate, "%content%", content, 1)
	replacedEmailSender := strings.Replace(replacedContent, "%admin_email%", from, 1)
	return replacedEmailSender

}

func setupMailServer() *mail.SMTPServer {

	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	return server

}

func setupMailClient(server *mail.SMTPServer) *mail.SMTPClient {

	client, err := server.Connect()
	if err != nil {
		AppConfig.ErrorLog.Println(err)
	}

	return client

}
