package spoonacular

type SpoonacularAPI struct {
	BaseURL string
	APIKey  string
}

func NewSpoonacularAPI(
	baseURL string,
	APIKey string,
) *SpoonacularAPI {
	return &SpoonacularAPI{
		BaseURL: baseURL,
		APIKey:  APIKey,
	}
}

type Ingredient struct {
	Results []*IngredientResult `json:"results"`
}

type IngredientResult struct {
	ID            int      `json:"id"`
	Original      string   `json:"original"`
	OriginalName  string   `json:"originalName"`
	Name          string   `json:"name"`
	Amount        float64  `json:"amount"`
	Unit          string   `json:"unit"`
	UnitShort     string   `json:"unitShort"`
	UnitLong      string   `json:"unitLong"`
	PossibleUnits []string `json:"possibleUnits"`
	EstimatedCost struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	} `json:"estimatedCost"`
	Consistency  string     `json:"consistency"`
	Aisle        string     `json:"aisle"`
	Image        string     `json:"image"`
	Meta         []struct{} `json:"meta"`
	Nutrition    Nutrition  `json:"nutrition"`
	CategoryPath []string   `json:"categoryPath"`
}

type Nutrient struct {
	Name                string  `json:"name"`
	Amount              float64 `json:"amount"`
	Unit                string  `json:"unit"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds"`
}

type Nutrition struct {
	Nutrients        []Nutrient `json:"nutrients"`
	Properties       []struct{} `json:"properties"`
	CaloricBreakdown struct {
		PercentProtein float64 `json:"percentProtein"`
		PercentFat     float64 `json:"percentFat"`
		PercentCarbs   float64 `json:"percentCarbs"`
	} `json:"caloricBreakdown"`
	WeightPerServing struct {
		Amount float64 `json:"amount"`
		Unit   string  `json:"unit"`
	} `json:"weightPerServing"`
}
