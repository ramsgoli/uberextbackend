package uber

import (
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
)

type request struct {
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
}

type UberResponseResult struct {
	DisplayName string `json:"localized_display_name"`
	Distance float32 `json:"distance"`
	Estimate string `json:"estimate"`
}

type ClientResponse struct {
    Message string `json:"message"`
    Prices []UberResponseResult `json:"prices"`
}

func GetUberEstimate(w http.ResponseWriter, r *http.Request) {
	if r.URL == nil {
		http.Error(w, "Please provide destination coordinates", 400)
		return
	}

	clientReq := r.URL.Query()
	if clientReq["end_latitude"] == nil || clientReq["end_longitude"] == nil {
		http.Error(w, "Please provide end_latitude and end_longitude", 400)
		return
	}


	UberReq, getError := http.NewRequest("GET", "https://api.uber.com/v1.2/estimates/price", nil)
	res := ClientResponse{}
	if getError != nil {
		res.Message = "Can not get an estimate right now"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	query := UberReq.URL.Query()
	query.Add("end_latitude", clientReq.Get("end_latitude"))
	query.Add("end_longitude", clientReq.Get("end_longitude"))
	UberReq.URL.RawQuery = query.Encode()

	// Set header
	var buffer bytes.Buffer
	buffer.WriteString("Bearer ")
	buffer.WriteString("vSzuZLd5Hxgs6RfSxD36n7ZHr0oHnP7euXbfb6g0")
	fmt.Println(buffer.String())
	UberReq.Header.Set("Authorization", buffer.String())

	client := http.Client{}
	uberResp, uberGetError := client.Do(UberReq)
	if uberGetError != nil {
		res.Message = "Can not get an estimate right now"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	json.NewDecoder(uberResp.Body).Decode(&res)
	respError := json.NewEncoder(w).Encode(res)
	if respError != nil {
		http.Error(w, respError.Error(), 500)
		return
	}
}
