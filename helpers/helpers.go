package helpers

import (
	"io/ioutil"
	"log"
	"net/http"
)

// FetchData opens urls and return the body of request
func FetchData(URL string) []byte {

	res, err := http.Get(URL)
	if err != nil {
		log.Print("Failed to make request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("Failed to read response")
	}

	return body

}
