package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// User is a struct representation of the user JSON object from tock
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// ReportingPeriod is a struct representation of the reporting_period JSON object from tock
type ReportingPeriod struct {
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	WorkingHours int    `json:"working_hours"`
}

// APIPages is a struct representation of a API page response from tock
type APIPages struct {
	Count   int    `json:"count"`
	NextURL string `json:"next"`
	PrevURL string `json:"previous"`
}

// ReportingPeriodAuditList is a struct representation of an API response from
//the Reporting Period Audit list endpoint
type ReportingPeriodAuditList struct {
	APIPages
	ReportingPeriods []ReportingPeriod `json:"results"`
}

// ReportingPeriodAuditDetails is a struct representation of an API response
//from the Reporting Period Audit details endpoint
type ReportingPeriodAuditDetails struct {
	APIPages
	Users []User `json:"results"`
}

// fetchDataProxy makes a GET request while adding the SignRequest opens url and
// return the body of request
func (bot *Bot) fetchDataProxy(URL string) []byte {

	req, err := http.NewRequest("GET", URL, nil)
	bot.Auth.SignRequest(req)
	if err != nil {
		log.Print("Failed to make request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print("Failed to make request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("Failed to read response")
	}

	return body
}

// Function for collecting the current reporting period
func (bot *Bot) fetchCurrentReportingPeriod() string {

	var data ReportingPeriodAuditList

	URL := fmt.Sprintf(os.Getenv("AUDIT_ENDPOINT"))

	body := bot.fetchDataProxy(URL)

	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}

	return data.ReportingPeriods[0].StartDate
}

// FetchTockUsers is a function for collecting all the users who have not
// filled out thier time sheet for the current period
func (bot *Bot) FetchTockUsers() *ReportingPeriodAuditDetails {

	var data ReportingPeriodAuditDetails
	timePeriod := bot.fetchCurrentReportingPeriod()

	URL := fmt.Sprintf("%s%s/", bot.AuditEndpoint, timePeriod)
	body := bot.fetchDataProxy(URL)

	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}
	return &data
}
