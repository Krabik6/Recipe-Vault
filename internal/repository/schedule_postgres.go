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
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	fillScheduleQuery := fmt.Sprintf(`INSERT INTO %s ("dateOf", "breakfastId", "lunchId", "dinnerId", "userId") values (date('%s'), $1, $2, $3, $4) RETURNING id`, scheduleTable, schedule.Date)
	row := tx.QueryRow(fillScheduleQuery, schedule.BreakfastId, schedule.LunchId, schedule.DinnerId, userId)

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()

}

func (s *SchedulePostgres) GetAllSchedule(userId int) ([]models.ScheduleOutput, error) {
	var output []models.ScheduleOutput

	getAllRecipesQuery :=
		fmt.Sprintf(`SELECT 
       r.title        as "BreakfastTitle",
       r.description  as "BreakfastDescription",
       r2.title       as "LunchTitle",
       r2.description as "LunchDescription",
       r3.title       as "DinnerTitle",
       r3.description as "DinnerDescription"
FROM %s as st
         JOIN recipes r on r.id = st."breakfastId"
         JOIN recipes r2 on r2.id = st."lunchId"
         JOIN recipes r3 on r3.id = st."dinnerId"
    	WHERE st."userId" = $1
`, scheduleTable)

	err := s.db.Select(&output, getAllRecipesQuery, userId)
	if err != nil {
		return output, err
	}

	return output, err
}
func (s *SchedulePostgres) GetScheduleByDate(userId int, date string) (models.ScheduleOutput, error) {
	var output models.ScheduleOutput

	GetScheduleByDateQuery :=
		fmt.Sprintf(`SELECT 
    r.title        as "BreakfastTitle",
    r.description  as "BreakfastDescription",
    r2.title       as "LunchTitle",
    r2.description as "LunchDescription",
    r3.title       as "DinnerTitle",
    r3.description as "DinnerDescription"
	FROM %s as st
	JOIN recipes r on r.id = st."breakfastId"
	JOIN recipes r2 on r2.id = st."lunchId"
	JOIN recipes r3 on r3.id = st."dinnerId"
	WHERE schedule."dateOf" = date('%s') and st."userId" = $1
`, scheduleTable, date)

	log.Println(GetScheduleByDateQuery)
	err := s.db.Get(&output, GetScheduleByDateQuery, userId)
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

	if input.BreakfastId != nil {
		setValues = append(setValues, fmt.Sprintf("breakfast_id=$%d", argId))
		args = append(args, *input.BreakfastId)
		argId++
	}

	if input.LunchId != nil {
		setValues = append(setValues, fmt.Sprintf("lunch_id=$%d", argId))
		args = append(args, *input.LunchId)
		argId++
	}

	if input.DinnerId != nil {
		setValues = append(setValues, fmt.Sprintf("dinner_id=$%d", argId))
		args = append(args, *input.DinnerId)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE schedule.date_of=date('%s')", scheduleTable, setQuery, date)
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
