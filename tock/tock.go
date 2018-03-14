package tockPackage

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/18F/angrytock/helpers"
	"github.com/cloudfoundry-community/go-cfenv"
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
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	ExactWorkingHours int    `json:"exact_working_hours"`
	MinWorkingHours   int    `json:"min_working_hours"`
	MaxWorkingHours   int    `json:"max_working_hours"`
}

// ReportingPeriodAuditList is a struct representation of an API response from
//the Reporting Period Audit list endpoint
type ReportingPeriodAuditList struct {
	ReportingPeriods []ReportingPeriod
}

// ReportingPeriodAuditDetails is a struct representation of an API response
//from the Reporting Period Audit details endpoint
type ReportingPeriodAuditDetails struct {
	Users []User
}

// Tock struct contains the audit endpoint and methods associated with Tock
type Tock struct {
	// Get Audit endpoint
	TockURL       string
	UserTockURL   string
	AuditEndpoint string
	DataFetcher   *helpers.DataFetcher
}

// InitTock initalizes the tock struct
func InitTock() *Tock {
	appEnv, _ := cfenv.Current()
	appService, _ := appEnv.Services.WithName("angrytock-credentials")

	// Get the tock url
	tockURL := fmt.Sprint(appService.Credentials["TOCK_URL"])
	if tockURL == "" {
		log.Fatal("TOCK_URL environment variable not found")
	}
	userTockURL := fmt.Sprint(appService.Credentials["USER_TOCK_URL"])
	if userTockURL == "" {
		log.Fatal("USER_TOCK_URL environment variable not found")
	}
	auditEndpoint := tockURL + "/api/reporting_period_audit"
	// Initalize a new data fetcher
	dataFetcher := helpers.NewDataFetcher(helpers.FetchData)
	return &Tock{tockURL, userTockURL, auditEndpoint, dataFetcher}
}

// fetchCurrentReportingPeriod gets the latest reporting time period that
// has happend
func fetchCurrentReportingPeriod(data *ReportingPeriodAuditList) string {
	currentPeriodIndex := 0
	for idx, period := range data.ReportingPeriods {
		endDate, _ := time.Parse("2006-01-02", period.EndDate)
		if endDate.Before(time.Now()) {
			currentPeriodIndex = idx
			break
		}
	}
	return data.ReportingPeriods[currentPeriodIndex].StartDate
}

// fetchReportingPeriod collects the current reporting period
func (tock *Tock) fetchReportingPeriod() string {
	var data ReportingPeriodAuditList
	URL := fmt.Sprintf("%s.json", tock.AuditEndpoint)
	body := tock.DataFetcher.FetchData(URL)
	err := json.Unmarshal(body, &data.ReportingPeriods)
	if err != nil {
		log.Print(err)
	}
	return fetchCurrentReportingPeriod(&data)
}

// FetchTockUsers is a function for collecting all the users who have not
// filled out thier time sheet for the current period
func (tock *Tock) FetchTockUsers(endpoint string) *ReportingPeriodAuditDetails {
	var data ReportingPeriodAuditDetails
	body := tock.DataFetcher.FetchData(endpoint)
	err := json.Unmarshal(body, &data.Users)
	if err != nil {
		log.Print(err)
	}
	return &data
}
