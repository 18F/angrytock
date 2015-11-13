package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type ReportingPeriod struct {
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	WorkingHours int    `json:"working_hours"`
}

type APIPages struct {
	Count   int    `json:"count"`
	NextUrl string `json:"next"`
	PrevUrl string `json:"previous"`
}

type ReportingPeriodAuditList struct {
	APIPages
	ReportingPeriods []ReportingPeriod `json:"results"`
}

type ReportingPeriodAuditDetails struct {
	APIPages
	Users []User `json:"results"`
}

func fetchCurrentReportingPeriod() string {

	var data ReportingPeriodAuditList

	url := fmt.Sprintf(os.Getenv("AUDIT_ENDPOINT"))
	res, err := http.Get(url)
	if err != nil {
		log.Print("Failed to make request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("Failed to read response")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}

	return data.ReportingPeriods[0].StartDate
}

func FetchTockUsers() *ReportingPeriodAuditDetails {

	var data ReportingPeriodAuditDetails
	timePeriod := fetchCurrentReportingPeriod()

	url := fmt.Sprintf("%s%s", os.Getenv("AUDIT_ENDPOINT"), timePeriod)
	res, err := http.Get(url)
	if err != nil {
		log.Print("Failed to make request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("Failed to read response")
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Print(err)
	}
	return &data
}
