package interfaces

import "github.com/ernesto27/docs/structs"

// crete interface
type DocDB interface {
	CreateDoc(doc structs.Doc) error
}
