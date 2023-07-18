package spoonacular

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (api *SpoonacularAPI) SearchIngredient(name string) (*Ingredient, error) {
	url := fmt.Sprintf("%s/food/ingredients/search?query=%s&apiKey=%s", api.BaseURL, name, api.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ingredient Ingredient
	err = json.Unmarshal(body, &ingredient)
	if err != nil {
		return nil, err
	}

	return &ingredient, nil
}
