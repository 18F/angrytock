package tockPackage

import (
	"testing"
	"time"

	"github.com/18F/tock-bot/helpers"
)

var test = []struct {
	Input  ReportingPeriodAuditList
	Output string
}{
	// Check with dates that already happend
	{
		ReportingPeriodAuditList{
			ReportingPeriods: []ReportingPeriod{
				ReportingPeriod{StartDate: "2014-01-07"},
				ReportingPeriod{StartDate: "2014-01-01"},
			},
		},
		"2014-01-07",
	},
	// Check with last date that hasn't occured
	{
		ReportingPeriodAuditList{
			ReportingPeriods: []ReportingPeriod{
				ReportingPeriod{StartDate: time.Now().Add(time.Hour * 48).Format("2006-01-02")},
				ReportingPeriod{StartDate: "2014-01-07"},
				ReportingPeriod{StartDate: "2014-01-01"},
			},
		},
		"2014-01-07",
	},
	// Check with last date that has occured
	{
		ReportingPeriodAuditList{
			ReportingPeriods: []ReportingPeriod{
				ReportingPeriod{StartDate: time.Now().Add(time.Hour * -48).Format("2006-01-02")},
				ReportingPeriod{StartDate: "2014-01-07"},
				ReportingPeriod{StartDate: "2014-01-01"},
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
