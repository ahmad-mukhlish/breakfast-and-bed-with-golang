package main

import (
	"github.com/ahmad-mukhlish/breakfast-and-bed-with-golang/internal/model"
	mail "github.com/xhit/go-simple-mail/v2"
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

	email := mail.NewMSG()
	email.SetFrom(mailData.From).AddTo(mailData.To).SetSubject(mailData.Subject)
	email.SetBody(mail.TextHTML, mailData.Content)

	err := email.Send(client)
	if err != nil {
		AppConfig.ErrorLog.Print(err)
	} else {
		AppConfig.InfoLog.Print("Email Sent")
	}

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
