package main

import (
	"github.com/userao/url-shortener/config"
	"github.com/userao/url-shortener/pkg/app"
	connector "github.com/userao/url-shortener/pkg/db-connector"
	"github.com/userao/url-shortener/pkg/server"
)

func main() {
	dbConnection := connector.NewConnection(config.DbName, config.DbUser, config.DbPassword, config.DbHost, config.DbPort)

	server := server.NewServer(config.Host, config.Port)

	app := app.NewApp(server, dbConnection)
	app.StartApp()
}
