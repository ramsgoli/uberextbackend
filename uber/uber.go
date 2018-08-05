package uber

import (
	"encoding/json"
	"net/http"
	"bytes"
	"github.com/ramsgoli/uberextbackend/keys"
)

type UberResponseResult struct {
	DisplayName string `json:"localized_display_name"`
	Distance float32 `json:"distance"`
	Estimate string `json:"estimate"`
}

type ClientResponse struct {
    Message string `json:"message"`
    Prices []UberResponseResult `json:"prices"`
}

func GetUberEstimate(w http.ResponseWriter, r *http.Request, keys *keys.Keys) {
	if r.URL == nil {
		http.Error(w, "Please provide destination coordinates", 400)
		return
	}

	clientReq := r.URL.Query()
	if clientReq["end_latitude"] == nil || clientReq["end_longitude"] == nil {
		http.Error(w, "Please provide end_latitude and end_longitude", 400)
		return
	}

	uberReq, getError := http.NewRequest("GET", "https://api.uber.com/v1.2/estimates/price", nil)
	res := ClientResponse{}
	if getError != nil {
		res.Message = "Can not get an estimate right now"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res)
		return
	}

	query := uberReq.URL.Query()
	query.Add("start_latitude", clientReq.Get("start_latitude"))
	query.Add("start_longitude", clientReq.Get("start_longitude"))
	query.Add("end_latitude", clientReq.Get("end_latitude"))
	query.Add("end_longitude", clientReq.Get("end_longitude"))
	uberReq.URL.RawQuery = query.Encode()

	// Set header
	var buffer bytes.Buffer
	buffer.WriteString("Token ")
	buffer.WriteString("vSzuZLd5Hxgs6RfSxD36n7ZHr0oHnP7euXbfb6g0")
	uberReq.Header.Set("Authorization", buffer.String())

	client := http.Client{}
	uberResp, uberGetError := client.Do(uberReq)
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
