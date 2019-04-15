package derp

import (
	"github.com/benpate/derp"
	"github.com/benpate/drivers/sendgrid"
	"github.com/davecgh/go-spew/spew"
)

// Reporter knows how to email a derp alert using SendGrid
type Reporter struct {
	apiKey      string
	fromAddress string
	toAddress   string
}

// NewReporter returns a fully populated Reporter that can be used by the derp system.
func NewReporter(apiKey string, fromAddress string, toAddress string) *Reporter {

	return &Reporter{
		apiKey:      apiKey,
		fromAddress: fromAddress,
		toAddress:   toAddress,
	}
}

// Report uses SendGrid to send an error alert to a specific email box.
func (reporter *Reporter) Report(err *derp.Error) {

	sp := spew.NewDefaultConfig()
	sp.ContinueOnMethod = true

	email := sendgrid.NewEmail()

	email.SetFromEmail(reporter.fromAddress, "DERP Sencder")
	email.AddToEmail(reporter.toAddress, "DERP Recipient")
	email.SetSubject("DERP Error: " + err.Message)
	email.AddBodyText(sp.Sdump(err))
	email.SendEmail(reporter.fromAddress)
}
