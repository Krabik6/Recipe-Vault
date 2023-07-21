package models

// IngredientInput используется для получения данных об ингредиентах от пользователя
type IngredientInput struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}
type Ingredients []Ingredient

type Ingredient struct {
	ID            int      `db:"id"`
	Name          string   `db:"name"`
	Price         float64  `db:"price"`
	Unit          string   `json:"unit"`
	UnitShort     string   `json:"unitShort"`
	UnitLong      string   `json:"unitLong"`
	PossibleUnits []string `json:"possibleUnits"`
	Protein       float64  `db:"protein"`
	Fat           float64  `db:"fat"`
	Carbs         float64  `db:"carbs"`
	Aisle         string   `db:"aisle"`
	Image         string   `db:"image"`
	CategoryPath  []string `db:"categoryPath"`
	Consistency   string   `db:"consistency"`
	ExternalID    int      `db:"external_id"`
	Amount        float64  `db:"amount" json:"amount"`
}

// IngredientOutput используется для отправки данных об ингредиентах пользователю
type IngredientOutput struct {
	ID            int      `db:"id" json:"id,omitempty"`
	Name          string   `db:"name" json:"name,omitempty"`
	Price         float64  `db:"price" json:"price,omitempty"`
	Unit          string   `db:"unit" json:"unit,omitempty"`
	UnitShort     string   `db:"unitShort" json:"unit_short,omitempty"`
	UnitLong      string   `db:"unitLong" json:"unit_long,omitempty"`
	PossibleUnits []string `db:"possible_units" json:"possible_units,omitempty"`
	Protein       float64  `db:"protein" json:"protein,omitempty"`
	Fat           float64  `db:"fat" json:"fat,omitempty"`
	Carbs         float64  `db:"carbs" json:"carbs,omitempty"`
	Aisle         string   `db:"aisle" json:"aisle,omitempty"`
	Image         string   `db:"image" json:"image,omitempty"`
	CategoryPath  []string `db:"categoryPath" json:"category_path,omitempty"`
	Consistency   string   `db:"consistency" json:"consistency,omitempty"`
	ExternalID    int      `db:"external_id" json:"external_id,omitempty"`
	Amount        float64  `db:"amount" json:"amount"`
}
