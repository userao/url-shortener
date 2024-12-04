package connector

import (
	"fmt"

	"github.com/userao/url-shortener/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IConnection interface {
	InitConnection()
	CreateUrl(string) (string, error)
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
	FullUrl    string
	Hash       string
	ClickCount uint `gorm:"not null;default:0"`
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

func (c Connection) CreateUrl(full string) (string, error) {
	var url Url

	result := c.db.Where("full_url = ?", full).First(&url)

	if result.RowsAffected != 0 {
		return url.Hash, nil
	}

	url = Url{FullUrl: full}
	c.db.Create(&url)
	hash := utils.GetHash(url.ID)
	url.Hash = hash
	c.db.Save(&url)

	return hash, nil
}

func (c Connection) GetUrl(shortenedUrl string) (Url, error) {
	var url Url
	result := c.db.Where("shortened_url = ?", shortenedUrl).First(&url)
	if result.RowsAffected == 0 {
		return Url{}, fmt.Errorf("record for %s not found", shortenedUrl)
	}

	return url, nil
}

func (c Connection) IncreaseClickCount(url *Url) {
	url.ClickCount = url.ClickCount + 1
	c.db.Save(url)
}
