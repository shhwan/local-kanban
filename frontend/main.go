package main

import (
	"local-kanban/frontend/client"
	"local-kanban/frontend/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "static")

	apiClient := client.NewBackendClient()
	boardHandler := handler.NewBoardHandler(apiClient)

	e.GET("/", boardHandler.HandleBoard)
	e.POST("/tasks", boardHandler.HandleCreateTask)
	e.PATCH("/tasks/:id/stage", boardHandler.HandleChangeStage)
	e.POST("/tasks/:id/logs", boardHandler.HandleAddWorkLog)
	e.POST("/tasks/:id/notes", boardHandler.HandleAddNote)
	e.DELETE("/tasks/:id", boardHandler.HandleDeleteTask)

	e.Logger.Fatal(e.Start(":3000"))
}
