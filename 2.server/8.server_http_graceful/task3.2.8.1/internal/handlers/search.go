package handlers

import (
	"authorization_jwt/internal/models"
	"encoding/json"
	"net/http"
)

// @Summary Поиск адреса
// @Description Принимает текстовый запрос и возвращает список адресов
// @Tags address
// @Accept json
// @Produce json
// @Param request body models.SearchRequest true "Запрос"
// @Success 200 {object} models.SearchResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/address/search [post]
// @Security BearerAuth
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var searchReq models.SearchRequest

	if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	client := createDaDataClient()

	addresses, err := client.SearchAddress(searchReq.Query)
	if err != nil {
		handleAPIError(w, err)
		return
	}

	searchResp := models.SearchResponse{Addresses: convertToPointerSlice(addresses)}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResp)
}
