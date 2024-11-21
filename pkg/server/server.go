package server

import (
	"fmt"
	"log"
	"net/http"

	connector "github.com/userao/url-shortener/pkg/db-connector"
)

type IServer interface {
	InitServer(connector.IConnection)
	ListenAndServe()
}

type Server struct {
	host string
	port string
}

var server Server

var dbConnection *connector.Connection

func NewServer(h, p string) *Server {
	server = Server{h, p}
	return &server
}

func GetCurrentServer() *Server {
	return &server
}

// тут обработчики эндпоинтов
func (s Server) InitServer(conn connector.IConnection) {
	dbConnection = conn.(*connector.Connection)
	http.HandleFunc("GET /urls/:id", getUrl)
	http.HandleFunc("GET /urls", getAllUrls)
	http.HandleFunc("POST /urls/create", createUrl)
	http.HandleFunc("GET /{hash}", redirectToOriginalUrl)
}

func (s Server) ListenAndServe() {
	fmt.Printf("Server started at %s:%s\n", s.host, s.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", s.host, s.port), nil))
}
