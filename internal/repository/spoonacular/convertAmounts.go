package spoonacular

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ConversionResult struct {
	SourceAmount float64 `json:"sourceAmount"`
	SourceUnit   string  `json:"sourceUnit"`
	TargetAmount float64 `json:"targetAmount"`
	TargetUnit   string  `json:"targetUnit"`
}

func (api *SpoonacularAPI) ConvertAmounts(ingredientName string, sourceAmount float64, sourceUnit string, targetUnit string) (*ConversionResult, error) {
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

	var conversionResult ConversionResult
	err = json.Unmarshal(body, &conversionResult)
	if err != nil {
		return nil, err
	}

	return &conversionResult, nil
}
