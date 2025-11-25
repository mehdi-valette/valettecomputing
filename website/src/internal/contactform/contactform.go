package contactform

import (
	"fmt"
	"net/http"
	"net/smtp"

	"valette.software/internal/config"
	"valette.software/internal/page"
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

	err = page.DisplayContactFormSuccess(res)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func sendEmail(form ContactForm) error {
	config := config.GetConfig()

	smtpData := config.GetSmtp()

	message := "From: " + smtpData.From + "\nTo: " + smtpData.To[0] + "\nSubject: valette.software - " + form.subject + "\n\nName: " + form.name + "\nCompany: " + form.company + "\nContact: " + form.contact + "\nSubject: " + form.subject + "\n\n" + form.message

	return smtp.SendMail(smtpData.Host+":"+smtpData.Port, config.GetSmtpAuth(), "mehdi.valette@gmail.com", []string{"mehdi.valette@gmail.com"}, []byte(message))
}
