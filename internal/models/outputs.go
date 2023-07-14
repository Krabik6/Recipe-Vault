package models

type ScheduleOutput struct {
	Id            int      `json:"id,omitempty" db:"id"`
	Name          string   `json:"name,omitempty" db:"name"`
	AtTime        string   `json:"at_time,omitempty" db:"at_time"`
	Title         string   `json:"title,omitempty" db:"title"`
	Description   string   `json:"description,omitempty" db:"description"`
	Public        bool     `json:"public,omitempty" db:"public"`
	Cost          float64  `json:"cost,omitempty" db:"cost"`
	TimeToPrepare int      `json:"timeToPrepare,omitempty" db:"timeToPrepare"`
	Healthy       int      `json:"healthy,omitempty" db:"healthy"`
	ImageURLs     []string `json:"imageURLs,omitempty" db:"imageURLs"`
}

/*

SELECT  m.id, m.name, m.at_time, r.title, r.description, r.public, r.cost, r."timeToPrepare", r.healthy from meal m
    JOIN mealrecipes mr on m.id = mr."mealId" JOIN recipes r on mr."recipeId" = r.id
        WHERE m.user_id = 14
          AND m.at_time >= '2023-01-17'
          AND m.at_time <= TIMESTAMP '2023-01-17' + INTERVAL '1 days';
*/
