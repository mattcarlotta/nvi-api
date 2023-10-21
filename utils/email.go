package utils

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type CustomEmail struct {
	M          *mail.SGMailV3
	P          *mail.Personalization
	Name       string
	Address    string
	TemplateID string
}

func (email *CustomEmail) setTemplateData(property string, value string) {
	email.P.SetDynamicTemplateData(property, value)
}

func (email *CustomEmail) send() error {
	email.M.SetTemplateID(email.TemplateID)
	email.M.SetFrom(mail.NewEmail("nvi", GetEnv("EMAIL_ADDRESS")))

	sendToEmailAddresses := []*mail.Email{mail.NewEmail(email.Name, email.Address)}
	email.P.AddTos(sendToEmailAddresses...)
	email.P.SetDynamicTemplateData("name", email.Name)
	email.P.SetDynamicTemplateData("unsubscribe", GetEnv("CLIENT_HOST")+"/settings/")
	email.P.SetDynamicTemplateData("unsubscribe_preferences", GetEnv("CLIENT_HOST")+"/settings/")

	email.M.AddPersonalizations(email.P)

	request := sendgrid.GetRequest(GetEnv("SEND_GRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(email.M)
	_, err := sendgrid.API(request)
	return err
}

func SendAccountVerificationEmail(name string, address string, token string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	email := CustomEmail{
		M:          mail.NewV3Mail(),
		P:          mail.NewPersonalization(),
		Name:       name,
		Address:    address,
		TemplateID: GetEnv("SEND_GRID_VERIFICATION_TEMPLATE_ID"),
	}

	email.setTemplateData("verify_link", GetEnv("CLIENT_HOST")+"/verify?token="+token)

	return email.send()
}

func SendPasswordResetEmail(name string, address string, token string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	email := CustomEmail{
		M:          mail.NewV3Mail(),
		P:          mail.NewPersonalization(),
		Name:       name,
		Address:    address,
		TemplateID: GetEnv("SEND_GRID_PASSWORD_RESET_TEMPLATE_ID"),
	}

	email.setTemplateData("reset_password_link", GetEnv("CLIENT_HOST")+"/reset-password?token="+token)

	return email.send()
}

func SendPasswordResetConfirmationEmail(name string, address string) error {
	if GetEnv("IN_TESTING") != "false" {
		return nil
	}

	email := CustomEmail{
		M:          mail.NewV3Mail(),
		P:          mail.NewPersonalization(),
		Name:       name,
		Address:    address,
		TemplateID: GetEnv("SEND_GRID_PASSWORD_RESET_CONFIRMATION_TEMPLATE_ID"),
	}

	email.setTemplateData("contact_us_link", GetEnv("CONTACT_US_LINK"))

	return email.send()
}
