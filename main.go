package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/twilio/twilio-go"
	controller "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/Controller"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/database"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/middlewares"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/service"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "db"
)

func main() {
	echoContext := echo.New()
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	m, err := migrate.New(
		"file://./migrations",
		psqlInfo)
	if err != nil {
		fmt.Println(err)
	}
	m.Down()
	if err := m.Up(); err != nil {
		fmt.Println(err)
	}

	appMiddleware := &middlewares.AppMiddleware{
		Logger: echoContext.Logger,
		DB:     db,
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: "AC7ed880bbbac683e5c3ff3b553631be20", //
		Password: "c049546a99d30d5e0db61ba98f7dff9a",   //
	})

	userDataLayer := database.NewUserDatalayerImpl(db)
	userService := service.NewUserServiceImpl(userDataLayer, client)
	controller.NewUserController(echoContext, userService, appMiddleware)

	aiDataLayer := database.NewAIDataLayerImpl(db)
	aiService := service.NewAIServiceImpl(aiDataLayer)
	controller.NewAIController(echoContext, aiService, appMiddleware)

	echoContext.Logger.Fatal(echoContext.Start(":1323"))
}
