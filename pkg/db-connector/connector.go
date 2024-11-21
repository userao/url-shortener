package connector

import (
	"fmt"

	"github.com/userao/url-shortener/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IConnection interface {
	InitConnection()
	CreateUrl(string, string) (string, error)
}

type Connection struct {
	dbName   string
	user     string
	password string
	host     string
	port     string
	db       *gorm.DB
}

type Url struct {
	gorm.Model
	FullUrl      string
	ShortenedUrl string
}

var connection Connection

func NewConnection(dbName, user, password, host, port string) *Connection {
	connection = Connection{dbName, user, password, host, port, nil}

	return &connection
}

func GetCurrentConnection() *Connection {
	return &connection
}

func (c *Connection) InitConnection() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", c.host, c.port, c.user, c.dbName, c.password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	c.db = db

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Url{})

	if err != nil {
		panic(err)
	}
}

func (c Connection) CreateUrl(full, host string) (string, error) {
	var url Url

	result := c.db.Where("full_url = ?", full).First(&url)

	if result.RowsAffected != 0 {
		return url.ShortenedUrl, nil
	}

	url = Url{FullUrl: full}
	c.db.Create(&url)
	hash := utils.GetHash(url.ID)
	shortenedUrl := hash
	url.ShortenedUrl = shortenedUrl
	c.db.Save(&url)

	return shortenedUrl, nil
}

func (c Connection) GetFullUrl(shortenedUrl string) (string, error) {
	var url Url
	result := c.db.Where("shortened_url = ?", shortenedUrl).First(&url)
	if result.RowsAffected == 0 {
		return "", fmt.Errorf("record for %s not found", shortenedUrl)
	}

	return url.FullUrl, nil
}
