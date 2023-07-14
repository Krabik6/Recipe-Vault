package repository

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"strings"
	"time"
)

type SchedulePostgres struct {
	db *sqlx.DB
}

func NewSchedulePostgres(db *sqlx.DB) *SchedulePostgres {
	return &SchedulePostgres{db: db}
}

//func (s *SchedulePostgres) FillSchedule(userId int, schedule models.Schedule) (int, error) {
//	tx, err := s.db.Begin()
//	if err != nil {
//		return 0, err
//	}
//
//	var id int
//	fillScheduleQuery := fmt.Sprintf(
//		`INSERT INTO %s ("date_of", "breakfast_id", "lunch_id", "dinner_id", "user_id") values (date('%s'), $1, $2, $3, $4) RETURNING id`,
//		scheduleTable, schedule.Date)
//	row := tx.QueryRow(fillScheduleQuery, schedule.BreakfastId, schedule.LunchId, schedule.DinnerId, userId)
//
//	if err := row.Scan(&id); err != nil {
//		tx.Rollback()
//		return 0, err
//	}
//
//	return id, tx.Commit()
//}

func (s *SchedulePostgres) CreateMeal(userId int, meal models.Meal) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	fillScheduleQuery := fmt.Sprintf(
		`INSERT INTO %s ("name", "at_time", "user_id") values ('%s', '%s', $1) RETURNING id`,
		mealTable, meal.Name, meal.AtTime)
	row := tx.QueryRow(fillScheduleQuery, userId)

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, value := range meal.Recipes {
		fillQuery := fmt.Sprintf(`
		INSERT INTO %s ("recipeId", "mealId") values ($1, $2)`,
			mealRecipesTable)
		_, err = tx.Exec(fillQuery, value, id)
		if err != nil {
			return 0, err
		}
	}

	return id, tx.Commit()

}
func (s *SchedulePostgres) GetAllSchedule(userID int) ([]models.ScheduleOutput, error) {
	var output []models.ScheduleOutput

	query := `
        SELECT
            m.id,
            m.name,
            m.at_time,
            r.title,
            r.description,
            r.public,
            r.cost,
            r."timeToPrepare",
            r.healthy,
            r."imageURLs"
        FROM
            meal m
        JOIN
            mealrecipes mr ON m.id = mr."mealId"
        JOIN
            recipes r ON mr."recipeId" = r.id
        WHERE
            m.user_id = $1
    `

	rows, err := s.db.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var schedule models.ScheduleOutput
		var imageURLs []string

		err := rows.Scan(
			&schedule.Id,
			&schedule.Name,
			&schedule.AtTime,
			&schedule.Title,
			&schedule.Description,
			&schedule.Public,
			&schedule.Cost,
			&schedule.TimeToPrepare,
			&schedule.Healthy,
			pq.Array(&imageURLs),
		)
		if err != nil {
			return nil, err
		}

		schedule.ImageURLs = imageURLs
		output = append(output, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func (s *SchedulePostgres) GetScheduleByPeriod(userID int, date string, dayPeriod int) ([]models.ScheduleOutput, error) {
	var output []models.ScheduleOutput
	query := `
        SELECT
            m.id,
            m.name,
            m.at_time,
            r.title,
            r.description,
            r.public,
            r.cost,
            r."timeToPrepare",
            r.healthy,
            r."imageURLs"
        FROM
            meal m
        JOIN
            mealrecipes mr ON m.id = mr."mealId"
        JOIN
            recipes r ON mr."recipeId" = r.id
        WHERE
            m.user_id = $1
            AND m.at_time >= $2
            AND m.at_time <= $2 + INTERVAL '%d days'
    `

	startDate, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return nil, err
	}

	query = fmt.Sprintf(query, dayPeriod)

	rows, err := s.db.Queryx(query, userID, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var schedule models.ScheduleOutput
		var imageURLs []string

		err := rows.Scan(
			&schedule.Id,
			&schedule.Name,
			&schedule.AtTime,
			&schedule.Title,
			&schedule.Description,
			&schedule.Public,
			&schedule.Cost,
			&schedule.TimeToPrepare,
			&schedule.Healthy,
			pq.Array(&imageURLs),
		)
		if err != nil {
			return nil, err
		}

		schedule.ImageURLs = imageURLs
		output = append(output, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
}
func (s *SchedulePostgres) UpdateSchedule(userID int, date string, input models.UpdateScheduleInput) error {
	db := s.db

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	parsedDate, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return err
	}

	if input.Recipes != nil {
		log.Println("i was also hereeee")

		// Get mealId for the given meal on the specified date
		getMealIdQuery := fmt.Sprintf(`SELECT "id" FROM %s WHERE "at_time" = $1 AND "user_id" = $2`, mealTable)
		var mealId int
		err := db.QueryRow(getMealIdQuery, parsedDate, userID).Scan(&mealId)
		if err != nil {
			log.Println("Error retrieving mealId:", err)
			return err
		}

		// Delete existing meal-recipe associations for the given mealId
		deleteAssociationsQuery := `DELETE FROM mealRecipes WHERE "mealId" = $1`
		_, err = db.Exec(deleteAssociationsQuery, mealId)
		if err != nil {
			log.Println("Error deleting meal-recipe associations:", err)
			return err
		}
		log.Println("Deleted existing meal-recipe associations")

		// Insert new meal-recipe associations for the given mealId
		for _, recipeID := range *input.Recipes {
			insertAssociationQuery := `INSERT INTO mealRecipes ("mealId", "recipeId") VALUES ($1, $2)`
			_, err = db.Exec(insertAssociationQuery, mealId, recipeID)
			if err != nil {
				log.Println("Error inserting meal-recipe association:", err)
				return err
			}
		}
		log.Println("Inserted new meal-recipe associations")
	}

	if input.AtTime != nil {
		setValues = append(setValues, fmt.Sprintf(`"at_time"=$%d`, argID))
		parsedInputDate, err := time.Parse("2006-01-02 15:04:05", *input.AtTime)
		if err != nil {
			return err
		}
		args = append(args, parsedInputDate)
		argID++

	}

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf(`"name"=$%d`, argID))
		args = append(args, *input.Name)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE "at_time"::date = $%d::date AND "user_id" = %d`, mealTable, setQuery, argID, userID)
	log.Println(query)
	args = append(args, parsedDate)

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}
	log.Println("i was hereeeeeeeeeee")

	return nil
}

func (s *SchedulePostgres) DeleteSchedule(userID int, date string) error {
	db := s.db

	// Parse the date string
	parsedDate, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return err
	}

	// Check if meal exists for the given user and date
	mealExistsQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE "at_time" = $1 AND "user_id" = $2`, mealTable)
	var count int
	err = db.Get(&count, mealExistsQuery, parsedDate, userID)
	if err != nil {
		return err
	}

	if count == 0 {
		// No meal scheduled for the given date
		return fmt.Errorf("no meal scheduled for the specified date")
	}

	// Delete the meal for the given user and date
	deleteScheduleQuery := fmt.Sprintf(`DELETE FROM %s WHERE "at_time" = $1 AND "user_id" = $2`, mealTable)
	_, err = db.Exec(deleteScheduleQuery, parsedDate, userID)
	if err != nil {
		return err
	}

	log.Println(deleteScheduleQuery)
	return nil
}

func (s *SchedulePostgres) DeleteMealsInRange(userID int, startDate string, endDate string) error {
	db := s.db

	// Check if meals exist for the given user and time range
	mealsExistQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE "at_time" >= DATE('%s') AND "at_time" <= DATE('%s') AND "user_id" = %d`, mealTable, startDate, endDate, userID)
	var count int
	err := db.Get(&count, mealsExistQuery)
	if err != nil {
		return err
	}

	if count == 0 {
		// No meals scheduled within the specified time range
		return fmt.Errorf("no meals scheduled within the specified time range")
	}

	// Delete the meals for the given user and time range
	deleteMealsQuery := fmt.Sprintf(`DELETE FROM %s WHERE "at_time" >= DATE('%s') AND "at_time" <= DATE('%s') AND "user_id" = %d`, mealTable, startDate, endDate, userID)
	_, err = db.Exec(deleteMealsQuery)
	if err != nil {
		return err
	}

	log.Println(deleteMealsQuery)
	return nil
}
