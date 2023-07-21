package main

import (
	"fmt"
	"github.com/Krabik6/meal-schedule/internal/apiserver"
	"github.com/Krabik6/meal-schedule/internal/cloudinaryutil"
	"github.com/Krabik6/meal-schedule/internal/configs"
	"github.com/Krabik6/meal-schedule/internal/handler"
	"github.com/Krabik6/meal-schedule/internal/models"
	"github.com/Krabik6/meal-schedule/internal/repository"
	"github.com/Krabik6/meal-schedule/internal/repository/spoonacular"
	"github.com/Krabik6/meal-schedule/internal/service"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

func main() {

	cfgPath, err := configs.ParseFlags("./config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := configs.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg.DB.Password)
	log.SetFlags(log.Ldate | log.Llongfile)

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatalf("db %e", err)
	}
	cloudinary := cloudinaryutil.NewCloudinaryClient(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
	options := &models.IngredientSearchOptions{
		SortDirection: "desc",
		Offset:        0,
		Number:        1,
	}
	spoonacularAPI := spoonacular.NewSpoonacularAPI(cfg.Spoonacular.BaseURL, cfg.Spoonacular.APIKey, options)
	repos := repository.NewRepository(db, spoonacularAPI)

	services := service.NewService(repos, cloudinary)
	handlers := handler.NewHandler(services)
	server := new(apiserver.Server)

	if err := server.Run(cfg.Server.Port, handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

	fmt.Println(strconv.Atoi("15"))

}
