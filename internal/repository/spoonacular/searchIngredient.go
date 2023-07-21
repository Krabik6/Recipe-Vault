package spoonacular

import (
	"encoding/json"
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (api *SpoonacularAPI) SearchIngredient(Query string) (models.IngredientSearchResponse, error) {
	options := &models.IngredientSearchOptions{
		SortDirection: "desc",
		Offset:        0,
		Number:        1,
	}

	api.Options = options

	queryParams := url.Values{}
	queryParams.Set("query", Query)
	//queryParams.Set("addChildren", strconv.FormatBool(api.Options.AddChildren))
	//queryParams.Set("minProteinPercent", strconv.Itoa(api.Options.MinProteinPercent))
	//queryParams.Set("maxProteinPercent", strconv.Itoa(api.Options.MaxProteinPercent))
	//queryParams.Set("minFatPercent", strconv.Itoa(api.Options.MinFatPercent))
	//queryParams.Set("maxFatPercent", strconv.Itoa(api.Options.MaxFatPercent))
	//queryParams.Set("minCarbsPercent", strconv.Itoa(api.Options.MinCarbsPercent))
	//queryParams.Set("maxCarbsPercent", strconv.Itoa(api.Options.MaxCarbsPercent))
	//queryParams.Set("metaInformation", strconv.FormatBool(api.Options.MetaInformation))
	//queryParams.Set("intolerances", api.Options.Intolerances)
	//queryParams.Set("sort", api.Options.Sort)
	queryParams.Set("sortDirection", api.Options.SortDirection)
	queryParams.Set("offset", strconv.Itoa(api.Options.Offset))
	queryParams.Set("number", strconv.Itoa(api.Options.Number))
	url := fmt.Sprintf("%s/food/ingredients/search?%s", api.BaseURL, queryParams.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.IngredientSearchResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-RapidAPI-Host", "spoonacular-recipe-food-nutrition-v1.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", api.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.IngredientSearchResponse{}, fmt.Errorf("failed to make request to Spoonacular API: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.IngredientSearchResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var ingredientResults models.IngredientSearchResponse

	err = json.Unmarshal(body, &ingredientResults)
	if err != nil {
		return models.IngredientSearchResponse{}, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return ingredientResults, nil
}
