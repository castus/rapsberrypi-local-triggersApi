package sunPositionAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	cacheFile = "../cache.json"
)

type SunPosition struct {
	Sunrise time.Time `json:"sunrise"`
	Sunset  time.Time `json:"sunset"`
}

type Response struct {
	Result SunPosition `json:"results"`
	Status string      `json:"status"`
}

type SunPositionAPI struct {
	SunPosition SunPosition
}

func NewSunPositionAPI() *SunPositionAPI {
	return &SunPositionAPI{}
}

func (s *SunPositionAPI) Get() SunPosition {
	if s.hasCacheHit() {
		log.Println("Sun position cache hit, serving data from cache")
		data, err := GetDataFromCache()
		if err != nil {
			panic(err)
		}
		return SunPosition{
			Sunset:  data.Sunset,
			Sunrise: data.Sunrise,
		}
	}

	log.Println("Sun position cache miss, serving data from web")
	URL := fmt.Sprintf("https://api.sunrise-sunset.org/json?lat=%s&lng=%s&date=today&formatted=0",
		os.Getenv("LATITUDE"),
		os.Getenv("LONGITUDE"))
	resp, err := http.Get(URL)

	if err != nil {
		fmt.Printf("Request Failed: %s", err)
		panic("Request Failed")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Reading body failed: %s", err)
		panic("Reading body failed")
	}
	res, err := getSunData(body)
	s.SunPosition = res.Result

	SaveDataToFile(res.Result)

	return res.Result
}

func SaveDataToFile(position SunPosition) {
	file, _ := json.MarshalIndent(position, "", " ")
	_ = os.WriteFile(cacheFile, file, 0644)
}

func GetDataFromCache() (*SunPosition, error) {
	file, readErr := os.ReadFile(cacheFile)
	if readErr != nil {
		return nil, readErr
	}

	var result *SunPosition
	err := json.Unmarshal([]byte(file), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getSunData(body []byte) (*Response, error) {
	var s = new(Response)
	err := json.Unmarshal(body, &s)
	if err != nil {
		log.Println("Error reading response from Sun position API")
		log.Println(err)
	}
	return s, err
}

func (s *SunPositionAPI) hasCacheHit() bool {
	now := time.Now()
	data, err := GetDataFromCache()
	if err != nil {
		return false
	}
	sunrise := data.Sunrise

	return now.Day() == sunrise.Day() && now.Month() == sunrise.Month() && now.Year() == sunrise.Year()
}
