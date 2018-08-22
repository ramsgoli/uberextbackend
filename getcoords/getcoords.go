package getcoords

import (
	"encoding/json"
	"net/http"
)

type GoogleResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Geometry struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

func GetLocation(address string) (error, GoogleResponse) {
	googleReq, getErr := http.NewRequest("GET", "https://maps.Google.com/maps/api/geocode/json", nil)
	if getErr != nil {
		return getErr, GoogleResponse{}
	}
	q := googleReq.URL.Query()
	q.Add("address", address)
	googleReq.URL.RawQuery = q.Encode()

	client := http.Client{}
	googleResp, getErr := client.Do(googleReq)
	if getErr != nil {
		return getErr, GoogleResponse{}
	}

	var reqRes GoogleResponse
	reqErr := json.NewDecoder(googleResp.Body).Decode(&reqRes)
	if reqErr != nil {
		return reqErr, GoogleResponse{}
	}
	return nil, reqRes
}

