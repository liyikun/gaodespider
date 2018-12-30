package insertdb

import (
	"database/sql"
	"fmt"
)

type ConnParam struct {
	host     string
	post     int
	user     string
	password string
	dbname   string
}

func ConnectDb(c ConnParam) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.host, c.post, c.user, c.password, c.dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
