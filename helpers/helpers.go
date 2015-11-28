package helpers

import (
	"io/ioutil"
	"log"
	"net/http"
)

// FetchData opens urls and return the body of request
func FetchData(URL string) []byte {
	// Get url
	res, err := http.Get(URL)
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
