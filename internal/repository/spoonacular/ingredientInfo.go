package spoonacular

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (api *SpoonacularAPI) GetIngredientInfo(id int) (*IngredientResult, error) {
	url := fmt.Sprintf("%s/food/ingredients/%d/information?apiKey=%s", api.BaseURL, id, api.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ingredientResult IngredientResult
	err = json.Unmarshal(body, &ingredientResult)
	if err != nil {
		return nil, err
	}

	return &ingredientResult, nil
}
