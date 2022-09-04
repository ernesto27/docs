package db

import (
	"database/sql"

	"github.com/ernesto27/docs/structs"
)

type Mock struct {
}

func (m *Mock) New() (error, *sql.DB) {
	return nil, nil
}

func (m *Mock) CreateDoc(doc structs.Doc) error {
	return nil
}

func (m *Mock) GetDocByID(ID int) (structs.Doc, error) {
	return structs.Doc{}, nil
}
