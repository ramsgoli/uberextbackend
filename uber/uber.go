package uber

import (
	"encoding/json"
	"net/http"
	"bytes"
	"strconv"
	"github.com/ramsgoli/uberextbackend/keys"
	"github.com/ramsgoli/uberextbackend/getcoords"
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

func floatToString(f float64) (string) {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func GetUberEstimate(w http.ResponseWriter, r *http.Request, keys *keys.Keys) {
	if r.URL == nil {
		http.Error(w, "Please provide destination coordinates", 400)
		return
	}

	clientReq := r.URL.Query()
	if clientReq["start_latitude"] == nil || clientReq["start_longitude"] == nil || clientReq["address"] == nil {
		http.Error(w, "Please provide start_latitude, start_longitude, and address", 400)
		return
	}

	getCoordsErr, coordinates := getcoords.GetLocation(clientReq.Get("address"))
	if getCoordsErr != nil {
		http.Error(w, "Could not fetch the coordinates of the address provided", 400)
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
	query.Add("end_latitude", floatToString(coordinates.Results[0].Geometry.Location.Lat))
	query.Add("end_longitude", floatToString(coordinates.Results[0].Geometry.Location.Lng))
	uberReq.URL.RawQuery = query.Encode()

	// Set header
	var buffer bytes.Buffer
	buffer.WriteString("Token ")
	buffer.WriteString(keys.UberKey)
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
