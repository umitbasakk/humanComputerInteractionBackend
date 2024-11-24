package main

import (
	"database/sql"
	"fmt"
	"os"

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

	if err := m.Up(); err != nil {
		fmt.Println(err)
	}

	appMiddleware := &middlewares.AppMiddleware{
		Logger: echoContext.Logger,
		DB:     db,
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("ACCOUNTSID"),
		Password: os.Getenv("AUTHTOKEN"),
	})

	datalayer := database.NewUserDatalayerImpl(db)
	userService := service.NewUserServiceImpl(datalayer, client)
	controller.NewUserController(echoContext, userService, appMiddleware)
	echoContext.Logger.Fatal(echoContext.Start(":1323"))
}
