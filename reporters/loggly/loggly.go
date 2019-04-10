package derp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Reporter knows how to report a derp error to the loggly.com service
type Reporter struct {
	url    string
	token string
}

// NewReporter returns a fully populated Reporter to be used by the derp system
func NewReporter(url string, token string) *Reporter {
	return &Reporter{
		url:    url,
		token: token,
	}
}

// Report sends a derp.Error to the loggly.com web service.
func (reporter *Reporter) Report(*derp.Error) {

		body, err := json.Marshal(result)

		if err != nil {
			fmt.Println("Loggly: error marshalling JSON: " + err.Error())
			return
		}

		url := Reporter.url + Reporter.token + "/tag/" + result.Location + "/"

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))

		if err != nil {
			fmt.Println("Loggly: error creating HTTP Request: " + err.Error())
			return
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(req)

		if err != nil {
			fmt.Println("Loggly: error sending report to Loggly: " + err.Error())
			return
		}

		if res.StatusCode != 200 {
			fmt.Println("Loggly: error response from Loggly service: ", res.StatusCode)
		}
	}
}
