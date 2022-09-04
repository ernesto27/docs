package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ernesto27/docs/structs"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	db sql.DB
}

func (m *Mysql) New() (error, *sql.DB) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))
	db, err := sql.Open("mysql", dataSourceName)
	defer db.Close()

	if err != nil {
		return err, nil
	}

	m.db = *db

	return nil, db
}

func (m *Mysql) CreateDoc(doc structs.Doc) error {
	query := "INSERT INTO docs (title, body) VALUES (?, ?)"
	_, err := m.db.Exec(query, doc.Title, doc.Body)
	return err
}
