package getcoords

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Lat  float32
	Long float32
}

type GoogleResponse struct {
	Results []Result
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
	fmt.Println(GoogleReq)
	resp, getErr := client.Do(GoogleReq)
	if getErr != nil {
		fmt.Println("get")
		http.Error(w, getErr.Error(), 500)
		return
	}

	var reqRes GoogleResponse
	reqErr := json.NewDecoder(resp.Body).Decode(&reqRes)
	if reqErr != nil {
		http.Error(w, reqErr.Error(), 500)
		return
	}
	fmt.Println(reqRes.Results[0].Geometry.Location.Lat)

	res := response{Lat: reqRes.Results[0].Geometry.Location.Lat, Long: reqRes.Results[0].Geometry.Location.Lng}
	respErr := json.NewEncoder(w).Encode(res)
	if respErr != nil {
		fmt.Println("res")
		http.Error(w, respErr.Error(), 500)
		return
	}
}
