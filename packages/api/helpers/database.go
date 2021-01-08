package helpers

import (
	"fmt"

	"github.com/jinzhu/gorm"
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
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", c.host, c.port, c.username, c.dbName, c.password,
	)

	if c.sslMode == "require" {
		connStr += fmt.Sprintf("&sslcert=%s&sslkey=%s&sslrootcert=%s", c.sslCert, c.sslKey, c.sslRootCert)
	}

	db, err := gorm.Open(c.db, connStr)

	return db, err
}
