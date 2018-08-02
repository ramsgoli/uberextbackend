package getcoords

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Lat  float32 `json:"lat"`
	Long float32 `json:"long"`
}

type GoogleResponse struct {
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

func GetLocation(w http.ResponseWriter, r *http.Request) {
	if r.URL == nil {
		http.Error(w, "Please provide a request body", 400)
		return
	}

	clientReq := r.URL.Query()
	if clientReq["address"] == nil {
		http.Error(w, "Please provide an address", 400)
		return
	}

	GoogleReq, getErr := http.NewRequest("GET", "https://maps.Google.com/maps/api/geocode/json", nil)
	if getErr != nil {
		http.Error(w, getErr.Error(), 500)
		return
	}
	q := GoogleReq.URL.Query()
	q.Add("address", clientReq.Get("address"))
	q.Add("key", "AIzaSyC5z_2ghHueFqRsB9-z9_5Vg1mvmBltL0A")
	GoogleReq.URL.RawQuery = q.Encode()

	client := http.Client{}
	resp, getErr := client.Do(GoogleReq)
	if getErr != nil {
		http.Error(w, getErr.Error(), 500)
		return
	}

	var reqRes GoogleResponse
	reqErr := json.NewDecoder(resp.Body).Decode(&reqRes)
	if reqErr != nil {
		http.Error(w, reqErr.Error(), 500)
		return
	}

	res := response{Lat: reqRes.Results[0].Geometry.Location.Lat, Long: reqRes.Results[0].Geometry.Location.Lng}
	respErr := json.NewEncoder(w).Encode(res)
	if respErr != nil {
		http.Error(w, respErr.Error(), 500)
		return
	}
}
