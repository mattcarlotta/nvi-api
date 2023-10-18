package utils

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendAccountVerificationEmail(name string, email string, token string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	m := mail.NewV3Mail()
	fromEmailAddress := mail.NewEmail("nvi", GetEnv("EMAIL_ADDRESS"))
	m.SetFrom(fromEmailAddress)

	m.SetTemplateID(GetEnv("SEND_GRID_VERIFICATION_TEMPLATE_ID"))
	p := mail.NewPersonalization()
	toEmailAddresses := []*mail.Email{
		mail.NewEmail(name, email),
	}
	p.AddTos(toEmailAddresses...)

	p.SetDynamicTemplateData("name", name)
	hostURL := GetEnv("CLIENT_HOST")
	verifyLink := hostURL + "/verify?token=" + token
	p.SetDynamicTemplateData("verify_link", verifyLink)
	unsubLink := hostURL + "/settings/"
	p.SetDynamicTemplateData("unsubscribe", unsubLink)
	p.SetDynamicTemplateData("unsubscribe_preferences", unsubLink)

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(GetEnv("SEND_GRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	_, err := sendgrid.API(request)
	return err
}
