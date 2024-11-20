package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	controller "github.com/umitbasakk/humanComputerInteractionBackend/UserStore/Controller"
	"github.com/umitbasakk/humanComputerInteractionBackend/UserStore/database"
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

	datalayer := database.NewUserDatalayerImpl(db)
	userService := service.NewUserServiceImpl(datalayer)
	controller.NewUserController(echoContext, userService)

	echoContext.Logger.Fatal(echoContext.Start(":1323"))
}
