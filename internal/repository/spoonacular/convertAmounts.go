package spoonacular

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"io/ioutil"
	"net/http"
)

func (api *SpoonacularAPI) ConvertAmounts(ingredientName string, sourceAmount float64, sourceUnit string, targetUnit string) (*models.ConversionResult, error) {
	url := fmt.Sprintf("%s/food/ingredients/convert?ingredientName=%s&sourceAmount=%f&sourceUnit=%s&targetUnit=%s&apiKey=%s", api.BaseURL, ingredientName, sourceAmount, sourceUnit, targetUnit, api.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var conversionResult models.ConversionResult
	err = json.Unmarshal(body, &conversionResult)
	if err != nil {
		return nil, err
	}

	return &conversionResult, nil
}
