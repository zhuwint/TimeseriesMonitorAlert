package mysql

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Account struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (m Account) Validate() error {
	if m.Username == "" {
		return fmt.Errorf("mysql: username could not be empty")
	}
	if m.Password == "" {
		return fmt.Errorf("mysql: password could not be empty")
	}
	if m.Address == "" {
		return fmt.Errorf("mysql: address could not be empty")
	}
	if m.Database == "" {
		return fmt.Errorf("mysql: database could not be empty")
	}
	return nil
}

type Connector struct {
	Address  string
	Database string
	Username string
	Password string
	DB       *gorm.DB
}

var (
	mysqlClient *Connector
	mysqlOnce   = &sync.Once{}
)

func newConnector(address, database, username, password string) (*Connector, error) {
	link := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, address, database)
	db, err := gorm.Open(mysql.Open(link), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		return nil, err
	}

	return &Connector{
		Address:  address,
		Database: database,
		Username: username,
		Password: password,
		DB:       db,
	}, nil
}

func InitMysqlClient(c Account) {
	mysqlOnce.Do(func() {
		conn, err := newConnector(c.Address, c.Database, c.Username, c.Password)
		if err != nil {
			panic(fmt.Errorf("init mysql conn failed %s", err.Error()))
		}
		mysqlClient = conn
	})
}

func GetClient() *gorm.DB {
	if mysqlClient == nil {
		panic("mysql client not init yeat")
	}
	return mysqlClient.DB
}
