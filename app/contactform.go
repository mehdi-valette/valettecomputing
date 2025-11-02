package main

import (
	"fmt"
	"net/http"
	"net/smtp"
)

type ContactForm struct {
	name    string
	company string
	contact string
	subject string
	message string
}

func HandleContactFormRequest(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		fmt.Println(err)
		return
	}

	form := ContactForm{
		name:    req.Form.Get("name"),
		company: req.Form.Get("company"),
		contact: req.Form.Get("contact"),
		subject: req.Form.Get("subject"),
		message: req.Form.Get("message"),
	}

	err = sendEmail(form)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = templates.ExecuteTemplate(res, "contactformsuccess.html", nil)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func sendEmail(form ContactForm) error {
	config := GetConfig()

	smtpData := config.GetSmtp()

	message := "From: " + smtpData.from + "\nTo: " + smtpData.to[0] + "\nSubject: valette.software - " + form.subject + "\n\nName: " + form.name + "\nCompany: " + form.company + "\nContact: " + form.contact + "\nSubject: " + form.subject + "\n\n" + form.message

	return smtp.SendMail(smtpData.host+":"+smtpData.port, config.GetSmtpAuth(), "mehdi.valette@gmail.com", []string{"mehdi.valette@gmail.com"}, []byte(message))
}
