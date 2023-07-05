package repository

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
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

func (s *SchedulePostgres) FillSchedule(userId int, meal models.Meal) (int, error) {
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

//"id" serial primary key,
//"name" varchar, --завтрак
//"at_time" timestamp not null, -- 10.04.2045, 10:00 по мск
//"user_id" int references users (id) on delete cascade not null, -- я (userId 5323)
//constraint userId_dateOf_unique unique("user_id", "at_time")

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
            r.healthy
        FROM
            meal m
        JOIN
            mealrecipes mr ON m.id = mr."mealId"
        JOIN
            recipes r ON mr."recipeId" = r.id
        WHERE
            m.user_id = $1
    `

	err := s.db.Select(&output, query, userID)
	if err != nil {
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
            r.healthy
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

	startDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	endDate := startDate.AddDate(0, 0, dayPeriod)

	query = fmt.Sprintf(query, dayPeriod)

	err = s.db.Select(&output, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (s *SchedulePostgres) UpdateSchedule(userId int, date string, input models.UpdateScheduleInput) error {
	db := s.db

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Date != nil {
		setValues = append(setValues, fmt.Sprintf(`"date_of"=$%d`, argId))
		args = append(args, *input.Date)
		argId++
	}

	if input.BreakfastId != nil {
		setValues = append(setValues, fmt.Sprintf(`"breakfast_id"=$%d`, argId))
		args = append(args, *input.BreakfastId)
		argId++
	}

	if input.LunchId != nil {
		setValues = append(setValues, fmt.Sprintf(`"lunch_id"=$%d`, argId))
		args = append(args, *input.LunchId)
		argId++
	}

	if input.DinnerId != nil {
		setValues = append(setValues, fmt.Sprintf(`"dinner_id"=$%d`, argId))
		args = append(args, *input.DinnerId)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE schedule."date_of"=date('%s') and schedule."user_id"=%d`, scheduleTable, setQuery, date, userId)

	args = append(args)

	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	return err

}
func (s *SchedulePostgres) DeleteSchedule(userId int, date string) error {
	db := s.db

	fillScheduleQuery := fmt.Sprintf(`DELETE FROM %s WHERE "date_of"=date('%s') and "user_id"=%d`, scheduleTable, date, userId)
	_, err := db.Exec(fillScheduleQuery)
	return err
}
