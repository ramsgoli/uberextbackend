package getcoords

import (
	"encoding/json"
	"net/http"
	"github.com/ramsgoli/uberextbackend/keys"
)

type googleResponse struct {
	Lat  float32 `json:"lat"`
	Long float32 `json:"long"`
}

type GooglegoogleResponse struct {
	Results []Result `json:"results"`
}

type Result struct {
	Geometry struct {
		Location struct {
			Lat float32 `json:"lat"`
			Lng float32 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

func GetLocation(w http.ResponseWriter, r *http.Request, keys *keys.Keys) {
	if r.URL == nil {
		http.Error(w, "Please provide a request body", 400)
		return
	}

	clientReq := r.URL.Query()
	if clientReq["address"] == nil {
		http.Error(w, "Please provide an address", 400)
		return
	}

	googleReq, getErr := http.NewRequest("GET", "https://maps.Google.com/maps/api/geocode/json", nil)
	if getErr != nil {
		http.Error(w, getErr.Error(), 500)
		return
	}
	q := googleReq.URL.Query()
	q.Add("address", clientReq.Get("address"))
	q.Add("key", "")
	googleReq.URL.RawQuery = q.Encode()

	client := http.Client{}
	googleResp, getErr := client.Do(googleReq)
	if getErr != nil {
		http.Error(w, getErr.Error(), 500)
		return
	}

	var reqRes GooglegoogleResponse
	reqErr := json.NewDecoder(googleResp.Body).Decode(&reqRes)
	if reqErr != nil {
		http.Error(w, reqErr.Error(), 500)
		return
	}

	res := googleResponse{Lat: reqRes.Results[0].Geometry.Location.Lat, Long: reqRes.Results[0].Geometry.Location.Lng}
	googleRespErr := json.NewEncoder(w).Encode(res)
	if googleRespErr != nil {
		http.Error(w, googleRespErr.Error(), 500)
		return
	}
}
