package utils

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(p *mail.Personalization, templateID string, name string, email string) error {
	m := mail.NewV3Mail()
	m.SetTemplateID(templateID)

	fromEmailAddress := mail.NewEmail("nvi", GetEnv("EMAIL_ADDRESS"))
	m.SetFrom(fromEmailAddress)

	toEmailAddresses := []*mail.Email{mail.NewEmail(name, email)}
	p.AddTos(toEmailAddresses...)

	p.SetDynamicTemplateData("name", name)
	p.SetDynamicTemplateData("unsubscribe", GetEnv("CLIENT_HOST")+"/settings/")
	p.SetDynamicTemplateData("unsubscribe_preferences", GetEnv("CLIENT_HOST")+"/settings/")

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(GetEnv("SEND_GRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err := sendgrid.API(request)
	return err
}

func SendAccountVerificationEmail(name string, email string, token string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	templateID := GetEnv("SEND_GRID_VERIFICATION_TEMPLATE_ID")

	p := mail.NewPersonalization()
	p.SetDynamicTemplateData("verify_link", GetEnv("CLIENT_HOST")+"/verify?token="+token)

	return SendEmail(p, templateID, name, email)
}

func SendPasswordResetEmail(name string, email string, token string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	templateID := GetEnv("SEND_GRID_PASSWORD_RESET_TEMPLATE_ID")

	p := mail.NewPersonalization()
	p.SetDynamicTemplateData("reset_password_link", GetEnv("CLIENT_HOST")+"/reset-password?token="+token)

	return SendEmail(p, templateID, name, email)
}

func SendPasswordResetConfirmationEmail(name string, email string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	templateID := GetEnv("SEND_GRID_PASSWORD_RESET_CONFIRMATION_TEMPLATE_ID")

	p := mail.NewPersonalization()
	p.SetDynamicTemplateData("contact_us_link", GetEnv("CONTACT_US_LINK"))

	return SendEmail(p, templateID, name, email)
}
