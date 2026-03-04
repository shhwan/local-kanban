package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"local-kanban/backend/handler"
	"local-kanban/backend/repository"
	"local-kanban/backend/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "kanban")
	dbPassword := getEnv("DB_PASSWORD", "kanban")
	dbName := getEnv("DB_NAME", "kanban")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("データベースへの接続に失敗しました: %v", err)
	}

	// スキーマはdb/init.sqlで管理。AutoMigrateは使用しない。

	taskRepo := repository.NewTaskRepository(db)
	labelRepo := repository.NewLabelRepository(db)
	stageRepo := repository.NewStageRepository(db)

	taskService := service.NewTaskService(taskRepo, stageRepo)
	labelService := service.NewLabelService(labelRepo)
	stageService := service.NewStageService(stageRepo)

	taskHandler := handler.NewTaskHandler(taskService)
	labelHandler := handler.NewLabelHandler(labelService)
	stageHandler := handler.NewStageHandler(stageService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowHeaders: []string{
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	}))

	api := e.Group("/api")

	api.GET("/tasks", taskHandler.GetAll)
	api.GET("/tasks/:id", taskHandler.GetByID)
	api.POST("/tasks", taskHandler.Create)
	api.PUT("/tasks/:id", taskHandler.Update)
	api.PATCH("/tasks/:id/stage", taskHandler.ChangeStage)
	api.DELETE("/tasks/:id", taskHandler.Delete)
	api.POST("/tasks/:id/logs", taskHandler.AddWorkLog)
	api.POST("/tasks/:id/notes", taskHandler.AddNote)

	api.GET("/labels", labelHandler.GetAll)
	api.POST("/labels", labelHandler.Create)

	api.GET("/stages", stageHandler.GetAll)
	api.POST("/stages", stageHandler.Create)

	e.Logger.Fatal(e.Start(":8080"))
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
