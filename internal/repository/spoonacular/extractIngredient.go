package spoonacular

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"io/ioutil"
	"net/http"
)

func (api *SpoonacularAPI) ExtractIngredient(text string) (*models.ExtractedIngredient, error) {
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

	var extractedIngredient models.ExtractedIngredient
	err = json.Unmarshal(body, &extractedIngredient)
	if err != nil {
		return nil, err
	}

	return &extractedIngredient, nil
}
