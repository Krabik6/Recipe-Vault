package models

type ConversionResult struct {
	SourceAmount float64 `json:"sourceAmount"`
	SourceUnit   string  `json:"sourceUnit"`
	TargetAmount float64 `json:"targetAmount"`
	TargetUnit   string  `json:"targetUnit"`
}

type ExtractedIngredient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}

type IngredientSearchOptions struct {
	AddChildren       bool
	MinProteinPercent int
	MaxProteinPercent int
	MinFatPercent     int
	MaxFatPercent     int
	MinCarbsPercent   int
	MaxCarbsPercent   int
	MetaInformation   bool
	Intolerances      string
	Sort              string
	SortDirection     string
	Offset            int
	Number            int
}

type IngredientSearchResponse struct {
	Results []IngredientAPIResponse `json:"results"`
}

type IngredientAPIResponse struct {
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

type IngredientInfoOptions struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}
