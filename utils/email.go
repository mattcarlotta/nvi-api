package utils

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(m *mail.SGMailV3, p *mail.Personalization) error {
	fromEmailAddress := mail.NewEmail("nvi", GetEnv("EMAIL_ADDRESS"))
	m.SetFrom(fromEmailAddress)

	hostURL := GetEnv("CLIENT_HOST")
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

func SendAccountVerificationEmail(name string, email string, token string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	m := mail.NewV3Mail()
	m.SetTemplateID(GetEnv("SEND_GRID_VERIFICATION_TEMPLATE_ID"))

	p := mail.NewPersonalization()
	toEmailAddresses := []*mail.Email{
		mail.NewEmail(name, email),
	}
	p.AddTos(toEmailAddresses...)

	p.SetDynamicTemplateData("name", name)
	p.SetDynamicTemplateData("verify_link", GetEnv("CLIENT_HOST")+"/verify?token="+token)

	return SendEmail(m, p)
}

func SendPasswordResetEmail(name string, email string, token string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	m := mail.NewV3Mail()
	m.SetTemplateID(GetEnv("SEND_GRID_PASSWORD_RESET_TEMPLATE_ID"))

	p := mail.NewPersonalization()
	toEmailAddresses := []*mail.Email{
		mail.NewEmail(name, email),
	}
	p.AddTos(toEmailAddresses...)

	p.SetDynamicTemplateData("name", name)
	resetPasswordLink := GetEnv("CLIENT_HOST") + "/reset-password?token=" + token
	p.SetDynamicTemplateData("reset_password_link", resetPasswordLink)

	return SendEmail(m, p)
}

func SendPasswordResetConfirmationEmail(name string, email string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	m := mail.NewV3Mail()
	m.SetTemplateID(GetEnv("SEND_GRID_PASSWORD_RESET_CONFIRMATION_TEMPLATE_ID"))

	p := mail.NewPersonalization()
	toEmailAddresses := []*mail.Email{
		mail.NewEmail(name, email),
	}
	p.AddTos(toEmailAddresses...)

	p.SetDynamicTemplateData("name", name)
	contactUsLink := GetEnv("CONTACT_US_LINK")
	p.SetDynamicTemplateData("contact_us_link", contactUsLink)

	return SendEmail(m, p)
}
