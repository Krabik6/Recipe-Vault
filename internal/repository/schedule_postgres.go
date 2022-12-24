package repository

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type SchedulePostgres struct {
	db *sqlx.DB
}

func NewSchedulePostgres(db *sqlx.DB) *SchedulePostgres {
	return &SchedulePostgres{db: db}
}

func (s *SchedulePostgres) FillSchedule(userId int, schedule models.Schedule) (int, error) {
	db := s.db
	fillScheduleQuery := fmt.Sprintf("INSERT INTO %s (date_of, breakfast_id, lunch_id, dinner_id) values (date('%s'), $1, $2, $3) RETURNING id", scheduleTable, schedule.Date)
	row := db.QueryRow(fillScheduleQuery, schedule.Breakfast, schedule.Lunch, schedule.Dinner)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SchedulePostgres) GetAllSchedule(userId int) ([]models.ScheduleOutput, error) {
	var output []models.ScheduleOutput

	getAllRecipesQuery :=
		fmt.Sprintf(`SELECT schedule.id as "id",
       date_of        as "DateOf",
       breakfast_id   as "BreakfastId",
       lunch_id       as "LunchId",
       dinner_id      as "DinnerId",
       r.title        as "BreakfastTitle",
       r.description  as "BreakfastDescription",
       r2.title       as "LunchTitle",
       r2.description as "LunchDescription",
       r3.title       as "DinnerTitle",
       r3.description as "DinnerDescription"
FROM schedule
         JOIN recipes r on r.id = schedule.breakfast_id
         JOIN recipes r2 on r2.id = schedule.lunch_id
         JOIN recipes r3 on r3.id = schedule.dinner_id
`)

	err := s.db.Select(&output, getAllRecipesQuery)
	if err != nil {
		return output, err
	}

	return output, err
}
func (s *SchedulePostgres) GetScheduleByDate(userId int, date string) (models.ScheduleOutput, error) {
	var output models.ScheduleOutput

	GetScheduleByDateQuery :=
		fmt.Sprintf(`SELECT schedule.id as "id",
       date_of        as "DateOf",
       breakfast_id   as "BreakfastId",
       lunch_id       as "LunchId",
       dinner_id      as "DinnerId",
       r.title        as "BreakfastTitle",
       r.description  as "BreakfastDescription",
       r2.title       as "LunchTitle",
       r2.description as "LunchDescription",
       r3.title       as "DinnerTitle",
       r3.description as "DinnerDescription"
FROM schedule
         JOIN recipes r on r.id = schedule.breakfast_id
         JOIN recipes r2 on r2.id = schedule.lunch_id
         JOIN recipes r3 on r3.id = schedule.dinner_id
WHERE schedule.date_of = date('%s')
`, date)

	log.Println(GetScheduleByDateQuery)
	err := s.db.Get(&output, GetScheduleByDateQuery)
	if err != nil {
		return output, err
	}

	return output, err
}
func (s *SchedulePostgres) UpdateSchedule(userId int, date string, input models.UpdateScheduleInput) error {
	db := s.db

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Date != nil {
		setValues = append(setValues, fmt.Sprintf("date=$%d", argId))
		args = append(args, *input.Date)
		argId++
	}

	if input.Breakfast != nil {
		setValues = append(setValues, fmt.Sprintf("breakfast=$%d", argId))
		args = append(args, *input.Breakfast)
		argId++
	}

	if input.Lunch != nil {
		setValues = append(setValues, fmt.Sprintf("lunch=$%d", argId))
		args = append(args, *input.Lunch)
		argId++
	}

	if input.Dinner != nil {
		setValues = append(setValues, fmt.Sprintf("dinner=$%d", argId))
		args = append(args, *input.Dinner)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE schedule.date_of=date('%s')", recipeTable, setQuery, date)
	args = append(args)

	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	return err

}
func (s *SchedulePostgres) DeleteSchedule(userId int, date string) error {
	db := s.db

	fillScheduleQuery := fmt.Sprintf("DELETE FROM %s WHERE schedule.date_of=date('%s')", scheduleTable, date)
	_, err := db.Exec(fillScheduleQuery)
	return err
}
