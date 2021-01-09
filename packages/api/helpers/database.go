package helpers

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db          string
	host        string
	dbName      string
	username    string
	password    string
	port        string
	timezone    string
	sslMode     string
	sslCert     string
	sslKey      string
	sslRootCert string
}

func NewDatabase(db, username, password, host, port, dbName, timezone, sslMode, sslCert, sslKey, sslRootCert string) *Database {
	return &Database{
		db:          db,
		username:    username,
		password:    password,
		host:        host,
		port:        port,
		dbName:      dbName,
		timezone:    timezone,
		sslMode:     sslMode,
		sslCert:     sslCert,
		sslKey:      sslKey,
		sslRootCert: sslRootCert,
	}
}

// connect to DB
func (c *Database) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", c.host, c.username, c.password, c.dbName, c.port, c.sslMode, c.timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}
