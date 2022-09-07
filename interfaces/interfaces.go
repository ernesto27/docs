package interfaces

import "github.com/ernesto27/docs/structs"

// crete interface
type DocDB interface {
	CreateDoc(doc structs.Doc) (int, error)
	GetDocByID(ID int) (structs.Doc, error)
	UpdateDocBodyByID(ID int, body string) (int, error)
	UpdateDocTitleByID(ID int, title string) (int, error)
}
