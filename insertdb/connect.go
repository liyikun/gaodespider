package insertdb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type ConnParam struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func ConnectDb(c ConnParam) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.host, c.port, c.user, c.password, c.dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
