package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cloudfoundry-community/go-cfenv"
)

// FetchData opens urls and return the body of request
func FetchData(URL string) []byte {
	appEnv, _ := cfenv.Current()
	appService, _ := appEnv.Services.WithName("angrytock-credentials")

	apiAuthToken := fmt.Sprintf("Token %s", appService.Credentials["TOCK_API_TOKEN"])
	if apiAuthToken == "" {
		log.Fatal("TOCK_API_TOKEN environment variable not found")
	}

	client := &http.Client{}
	// Get url
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", apiAuthToken)
	res, err := client.Do(req)
	if err != nil {
		log.Print("Failed to make request")
	}
	defer res.Body.Close()
	// Read body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("Failed to read response")
	}

	return body

}

// GenericDataFetcher is a generic function that takes a url string and returns
// a bytes
type GenericDataFetcher func(url string) []byte

// DataFetcher is a struct that holds a function that of type DataFetcher
type DataFetcher struct {
	GenericDataFetcherHolder GenericDataFetcher
}

// NewDataFetcher get
func NewDataFetcher(dataFetcher GenericDataFetcher) *DataFetcher {
	return &DataFetcher{GenericDataFetcherHolder: dataFetcher}
}

// FetchData get data
func (d *DataFetcher) FetchData(URL string) []byte {
	return d.GenericDataFetcherHolder(URL)
}
