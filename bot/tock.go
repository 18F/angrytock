package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Struct representation of the user JSON object from tock
type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Struct representation of the reporting_period JSON object from tock
type ReportingPeriod struct {
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	WorkingHours int    `json:"working_hours"`
}

// Struct representation of a API page response from tock
type APIPages struct {
	Count   int    `json:"count"`
	NextUrl string `json:"next"`
	PrevUrl string `json:"previous"`
}

// Struct representation of an API response from the Reporting Period Audit
// list endpoint
type ReportingPeriodAuditList struct {
	APIPages
	ReportingPeriods []ReportingPeriod `json:"results"`
}

// Struct representation of an API response from the Reporting Period Audit
// details endpoint
type ReportingPeriodAuditDetails struct {
	APIPages
	Users []User `json:"results"`
}

// Function for collecting the current reporting period
func fetchCurrentReportingPeriod() string {

	var data ReportingPeriodAuditList

	Url := fmt.Sprintf(os.Getenv("AUDIT_ENDPOINT"))

	body := fetchData(Url)

	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}

	return data.ReportingPeriods[0].StartDate
}

// Function for collecting all the users who have not filled out thier time sheet
// for the current period
func FetchTockUsers() *ReportingPeriodAuditDetails {

	var data ReportingPeriodAuditDetails
	timePeriod := fetchCurrentReportingPeriod()

	Url := fmt.Sprintf("%s%s", os.Getenv("AUDIT_ENDPOINT"), timePeriod)

	body := fetchData(Url)

	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}
	return &data
}
