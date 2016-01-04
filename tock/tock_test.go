package tockPackage

import (
	"testing"
	"time"

	"github.com/18F/angrytock/helpers"
)

var test = []struct {
	Input  ReportingPeriodAuditList
	Output string
}{
	// Check with dates that already happend
	{
		ReportingPeriodAuditList{
			ReportingPeriods: []ReportingPeriod{
				ReportingPeriod{StartDate: "2014-01-07", EndDate: "2014-01-12"},
				ReportingPeriod{StartDate: "2014-01-01", EndDate: "2014-01-05"},
			},
		},
		"2014-01-07",
	},
	// Check with last date that hasn't occured
	{
		ReportingPeriodAuditList{
			ReportingPeriods: []ReportingPeriod{
				ReportingPeriod{
					StartDate: time.Now().Add(time.Hour * 24 * 2).Format("2006-01-02"),
					EndDate:   time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02"),
				},
				ReportingPeriod{StartDate: "2014-01-07", EndDate: "2014-01-12"},
				ReportingPeriod{StartDate: "2014-01-01", EndDate: "2014-01-05"},
			},
		},
		"2014-01-07",
	},
	// Check with last date that has occured
	{
		ReportingPeriodAuditList{
			ReportingPeriods: []ReportingPeriod{
				ReportingPeriod{
					StartDate: time.Now().Add(time.Hour * -24 * 2).Format("2006-01-02"),
					EndDate:   time.Now().Add(time.Hour * -24 * 7).Format("2006-01-02"),
				},
				ReportingPeriod{StartDate: "2014-01-07", EndDate: "2014-01-12"},
				ReportingPeriod{StartDate: "2014-01-01", EndDate: "2014-01-05"},
			},
		},
		time.Now().Add(time.Hour * -48).Format("2006-01-02"),
	},
}

// Check that the most recent reporting period is returned
// without being a reporting period in the future
func TestFetchCurrentReportingPeriod(t *testing.T) {
	for _, test := range test {
		currentPeriod := fetchCurrentReportingPeriod(&test.Input)
		if currentPeriod != test.Output {
			t.Error(currentPeriod)
		}
	}
}

func mockDataFetcher(url string) []byte {
	if url == "AuditEndpoint" {
		return []byte(`{
			"count":62,
			"next":null,
			"previous":null,
			"results":[
				{"start_date":"2014-11-22","end_date":"2014-11-28","working_hours":40},
				{"start_date":"2014-11-15","end_date":"2014-11-21","working_hours":40}]
			}`)
	} else {
		return []byte(`{
	    "count":2,
	    "next":null,
	    "previous":null,
	    "results":[
	      {
	        "id":1,
	        "username":"user.one",
	        "first_name":"user",
	        "last_name":"one",
	        "email":"user.one@gsa.gov"},
	      {
	        "id":2,
	        "username":"user.two",
	        "first_name":"user",
	        "last_name":"two",
	        "email":"user.two@gsa.gov"
	      }
	    ]
	  }`)
	}
}

var tock = Tock{
	"TockURL",
	"UserTockURL",
	"AuditEndpoint",
	helpers.NewDataFetcher(mockDataFetcher),
}

func TestFetchTockReportingPeriods(t *testing.T) {
	reportingPeriod := tock.fetchReportingPeriod()
	if reportingPeriod != "2014-11-22" {
		t.Errorf(reportingPeriod)
	}
}

func TestFetchTockUsers(t *testing.T) {
	userData := tock.FetchTockUsers()
	if len(userData.Users) != 2 {
		t.Error(userData)
	}
}
