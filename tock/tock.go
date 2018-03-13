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
	ReportingPeriods []ReportingPeriod `json:"[]"`
}

// ReportingPeriodAuditDetails is a struct representation of an API response
//from the Reporting Period Audit details endpoint
type ReportingPeriodAuditDetails struct {
	APIPages
	Users []User `json:"results"`
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

// fetchCurrentReportingPeriod collects the current reporting period
func (tock *Tock) fetchReportingPeriod() string {
	var data ReportingPeriodAuditList
	URL := fmt.Sprintf("%s.json", tock.AuditEndpoint)
	body := tock.DataFetcher.FetchData(URL)
	err := json.Unmarshal(body, &data)
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
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}
	return &data
}

// TockUserGen returns a generator that returns a steram
// of user data by paging through the api
func (tock *Tock) TockUserGen() func() *ReportingPeriodAuditDetails {
	timePeriod := tock.fetchReportingPeriod()
	baseEndpoint := fmt.Sprintf("%s/%s.json", tock.AuditEndpoint, timePeriod)
	currentPage := 1
	newEndpoint := baseEndpoint + fmt.Sprintf("?page=%d", currentPage)
	return func() *ReportingPeriodAuditDetails {
		usersResponse := tock.FetchTockUsers(newEndpoint)
		currentPage++
		newEndpoint = baseEndpoint + fmt.Sprintf("?page=%d", currentPage)
		return usersResponse
	}
}

// UserApplier loops through users and applies a anonymous function to a list
// of late tock users
func (tock *Tock) UserApplier(applyFunc func(user User)) {
	// user Generator
	userGen := tock.TockUserGen()
	// get event indefinitely
	for {
		apiResponse := userGen()
		for _, user := range apiResponse.Users {
			applyFunc(user)
		}
		// Break loop if there are no more urls
		if apiResponse.NextURL == "" {
			break
		}
	}

}
