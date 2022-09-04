package db

import (
	"database/sql"
	"fmt"

	"github.com/ernesto27/docs/structs"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	db sql.DB
}

func (m *Mysql) New(user string, password string, host string, port string, name string) (error, *sql.DB) {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user,
		password, host, port, name)
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

func (m *Mysql) GetDocByID(ID int) (structs.Doc, error) {
	query := "SELECT id, title, body FROM docs WHERE id = ?"
	row := m.db.QueryRow(query, ID)
	doc := structs.Doc{}
	err := row.Scan(&doc.ID, &doc.Title, &doc.Body)
	return doc, err
}

func (m *Mysql) UpdateDocByID(ID int, body string) (int, error) {
	query := "UPDATE docs SET body = ? WHERE id = ?"
	res, err := m.db.Exec(query, body, ID)

	rowsAffected, errRows := res.RowsAffected()
	if errRows != nil {
		return 0, errRows
	}

	return int(rowsAffected), err
}
