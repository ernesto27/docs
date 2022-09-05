package db

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/ernesto27/docs/structs"
	"github.com/joho/godotenv"
)

var myDb Mysql

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	myDb = Mysql{}
	myDb.New(os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))

}

func TestGetDocById(t *testing.T) {

	type args struct {
		ID    int
		query string
	}

	type result struct {
		doc structs.Doc
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected result
	}{
		{
			name: "Get empty doc if ID is doc does not exist",
			args: args{
				ID: -9999,
			},
			expected: result{
				doc: structs.Doc{},
				err: errors.New("sql: no rows in result set"),
			},
		},

		{
			name: "Get result if ID doc exists on DB",
			args: args{
				ID:    1,
				query: "INSERT INTO docs (title, body) VALUES ('title1', 'body1')",
			},
			expected: result{
				doc: structs.Doc{
					ID:    1,
					Title: "title1",
					Body:  "body1",
				},
				err: errors.New("sql: no rows in result set"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myDb.db.Exec("TRUNCATE TABLE docs")
			if tt.args.query != "" {
				myDb.db.Exec(tt.args.query)
			}

			doc, err := myDb.GetDocByID(tt.args.ID)

			if doc != tt.expected.doc {
				t.Errorf("GetDocByID() = %v, want %v", doc, tt.expected.doc)
			}

			if err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("GetDocByID() = %v, want %v", err, tt.expected.err)
				}
			}
		})
	}

}

func TestUpdateDocById(t *testing.T) {

	type args struct {
		ID    int
		Body  string
		query string
	}

	type result struct {
		rows int
		err  error
	}

	tests := []struct {
		name     string
		args     args
		expected result
	}{
		{
			name: "Get rows affected = 0 if ID is doc does not exist",
			args: args{
				ID: -9999,
			},
			expected: result{
				rows: 0,
			},
		},
		{
			name: "Get rows affected = 1 if ID is doc exists on DB",
			args: args{
				ID:    1,
				Body:  "body1 updated",
				query: "INSERT INTO docs (title, body) VALUES ('title1', 'body1')",
			},
			expected: result{
				rows: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myDb.db.Exec("TRUNCATE TABLE docs")
			if tt.args.query != "" {
				myDb.db.Exec(tt.args.query)
			}

			rows, err := myDb.UpdateDocByID(tt.args.ID, tt.args.Body)

			if rows != tt.expected.rows {
				t.Errorf("UpdateDocByID() = %v, want %v", rows, tt.expected.rows)
			}

			if err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("UpdateDocByID() = %v, want %v", err, tt.expected.err)
				}
			}
		})
	}

}