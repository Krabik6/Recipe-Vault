package spoonacular

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ExtractedIngredient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

func (api *SpoonacularAPI) ExtractIngredient(text string) (*ExtractedIngredient, error) {
	url := fmt.Sprintf("%s/food/ingredients/extract?text=%s&apiKey=%s", api.BaseURL, text, api.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var extractedIngredient ExtractedIngredient
	err = json.Unmarshal(body, &extractedIngredient)
	if err != nil {
		return nil, err
	}

	return &extractedIngredient, nil
}
