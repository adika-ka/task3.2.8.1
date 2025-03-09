package handlers

import (
	"authorization_jwt/internal/models"
	"encoding/json"
	"net/http"
)

// @Summary Поиск адреса
// @Description Принимает текстовый запрос с координатами ("lat" - широта, "lng" - долгота) и возвращает список адресов
// @Tags address
// @Accept json
// @Produce json
// @Param request body models.GeocodeRequest true "Запрос"
// @Success 200 {object} models.GeocodeResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/address/geocode [post]
// @Security BearerAuth
func GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	var geocodeReq models.GeocodeRequest

	if err := json.NewDecoder(r.Body).Decode(&geocodeReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	client := createDaDataClient()

	addresses, err := client.GeocodeAddress(geocodeReq.Lat, geocodeReq.Lng)
	if err != nil {
		handleAPIError(w, err)
		return
	}

	geocodeResp := models.GeocodeResponse{Addresses: convertToPointerSlice(addresses)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(geocodeResp)
}
