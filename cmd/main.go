package main

import (
	"github.com/Krabik6/meal-schedule/internal/apiserver"
	"github.com/Krabik6/meal-schedule/internal/handler"
	"github.com/Krabik6/meal-schedule/internal/repository"
	"github.com/Krabik6/meal-schedule/internal/service"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "db",
		Port:     "5432",
		Username: "postgres",
		Password: "qwerty",
		DBName:   "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("db %e", err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(apiserver.Server)

	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

	//err = repository.CreateRecipe(models.Recipe{
	//	Id:          0,
	//	Title:        "borsh",
	//	Description: "russian soup",
	//}, db)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//_, err = repository.GetRecipeById(db)
	//if err != nil {
	//	panic(err)
	//}
}
