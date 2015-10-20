package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	OVER_QUERY_LIMIT = "OVER_QUERY_LIMIT"
	OK               = "OK"
)

var (
	customRandSource = rand.NewSource(time.Now().UnixNano())
	customRand       = rand.New(customRandSource)
)

type GoogleResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	Results      string `json:"results"`
}

func (this *GoogleResponse) String() string {
	if this.Status == OK {
		return fmt.Sprintf("%s", this.Status)
	}

	return fmt.Sprintf("ERROR %s %s %s", this.Status, this.ErrorMessage, this.Results)
}

func randomLatLon() (float32, float32) {
	nagative := 1
	if rand.Intn(2) == 0 {
		nagative = -1
	}
	lat := customRand.Float32() * 90 * float32(nagative)

	nagative = 1
	if rand.Intn(2) == 0 {
		nagative = -1
	}
	lon := customRand.Float32() * 180 * float32(nagative)

	return lat, lon
}

func doRequest(n int, lat, lon float32, requestUrl string) (*GoogleResponse, error) {
	res, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var data GoogleResponse
	err = json.Unmarshal(body, &data)

	return &data, nil
}

func printResAndStatistic(statistic map[string]int, requestUrl string, googleRes *GoogleResponse, n int) error {
	log.Println(fmt.Sprintf("\n%v: %s\n%v\n", n, requestUrl, googleRes))

	if n%10 == 0 {
		log.Println(fmt.Sprintf("\n%v\n=========\n", statistic))
	}

	if googleRes.Status == OVER_QUERY_LIMIT {
		return errors.New(googleRes.Status)
	}

	return nil
}

func main() {
	statistic := make(map[string]int)

	for i := 1; ; i++ {
		lat, lon := randomLatLon()
		requestUrl := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%v,%v", lat, lon)

		googleRes, err := doRequest(i, lat, lon, requestUrl)
		if err != nil {
			panic(err)
		}

		statistic[googleRes.Status]++

		err = printResAndStatistic(statistic, requestUrl, googleRes, i)
		if err != nil {
			panic(err)
		}

		time.Sleep(200 * time.Millisecond)
	}
}
