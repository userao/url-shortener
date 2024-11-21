package app

import (
	connector "github.com/userao/url-shortener/pkg/db-connector"
	"github.com/userao/url-shortener/pkg/server"
)

type App struct {
	server       server.IServer
	dbConnection connector.IConnection
}

var app App

func NewApp(s server.IServer, c connector.IConnection) *App {
	app = App{s, c}
	return &app
}

func (app *App) StartApp() {
	app.dbConnection.InitConnection()
	app.server.InitServer(app.dbConnection)
	app.server.ListenAndServe()
}
