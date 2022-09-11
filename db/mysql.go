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

func New(user string, password string, host string, port string, name string, tls string) (Mysql, error) {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", user,
		password, host, port, name, tls)
	db, err := sql.Open("mysql", dataSourceName)
	defer db.Close()

	if err != nil {
		return Mysql{}, err
	}

	return Mysql{
		db: *db,
	}, nil
}

func (m *Mysql) CreateDoc(doc structs.Doc) (int, error) {
	tx, errTransaction := m.db.Begin()
	if errTransaction != nil {
		return 0, errTransaction
	}

	query := "INSERT INTO docs (title, body) VALUES (?, ?)"
	res, err := tx.Exec(query, doc.Title, doc.Body)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	id, errID := res.LastInsertId()
	if errID != nil {
		tx.Rollback()
		return 0, errID
	}

	tx.Commit()
	return int(id), err
}

func (m *Mysql) GetDocByID(ID int) (structs.Doc, error) {
	query := "SELECT id, title, body FROM docs WHERE id = ?"
	row := m.db.QueryRow(query, ID)
	doc := structs.Doc{}
	err := row.Scan(&doc.ID, &doc.Title, &doc.Body)
	return doc, err
}

func (m *Mysql) UpdateDocBodyByID(ID int, body string) (int, error) {
	query := "UPDATE docs SET body = ? WHERE id = ?"
	res, err := m.db.Exec(query, body, ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, errRows := res.RowsAffected()
	if errRows != nil {
		return 0, errRows
	}

	return int(rowsAffected), nil
}

func (m *Mysql) UpdateDocTitleByID(ID int, title string) (int, error) {
	query := "UPDATE docs SET title = ? WHERE id = ?"
	res, err := m.db.Exec(query, title, ID)

	if err != nil {
		return 0, err
	}

	rowsAffected, errRows := res.RowsAffected()
	if errRows != nil {
		return 0, errRows
	}

	return int(rowsAffected), nil
}
