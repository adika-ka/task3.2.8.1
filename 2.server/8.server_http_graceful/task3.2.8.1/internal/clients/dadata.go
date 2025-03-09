package clients

import (
	"authorization_jwt/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DaDataClient struct {
	ApiKey     string
	HttpClient http.Client
}

type DaDataRequest struct {
	Query string `json:"query"`
}

type DaDataResponse struct {
	Suggestions []struct {
		Value string `json:"value"`
		Data  struct {
			Country string `json:"country"`
			City    string `json:"city"`
			Street  string `json:"street"`
			GeoLat  string `json:"geo_lat"`
			GeoLon  string `json:"geo_lon"`
		} `json:"data"`
	} `json:"suggestions"`
}

func (c *DaDataClient) makeDaDataRequest(apiURL string, payload interface{}) (*DaDataResponse, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("JSON encoding error: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Token "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var result DaDataResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}

func parseDaDataResponse(result *DaDataResponse) ([]models.Address, error) {
	if len(result.Suggestions) == 0 {
		return nil, fmt.Errorf("no addresses found")
	}

	addresses := []models.Address{}
	for _, suggestion := range result.Suggestions {
		address := models.Address{
			Street:  suggestion.Data.Street,
			City:    suggestion.Data.City,
			Country: suggestion.Data.Country,
			Lat:     suggestion.Data.GeoLat,
			Lng:     suggestion.Data.GeoLon,
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (c *DaDataClient) SearchAddress(query string) ([]models.Address, error) {
	result, err := c.makeDaDataRequest("https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address", DaDataRequest{Query: query})
	if err != nil {
		return nil, err
	}
	return parseDaDataResponse(result)
}

func (c *DaDataClient) GeocodeAddress(lat, lng string) ([]models.Address, error) {
	result, err := c.makeDaDataRequest("https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", map[string]string{
		"lat": lat,
		"lon": lng,
	})
	if err != nil {
		return nil, err
	}
	return parseDaDataResponse(result)
}
