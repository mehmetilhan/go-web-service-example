package storage

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

var DocumentDirectory = filepath.Join("uploads/documents")

type Document struct {
	FileName   string    `json:"name"`
	UploadDate time.Time `json:"uploadDate"`
}

func getDocuments() ([]Document, error) {
	documents := make([]Document, 0)

	files, err := ioutil.ReadDir(DocumentDirectory)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		documents = append(documents, Document{FileName: f.Name(), UploadDate: f.ModTime()})
	}

	return documents, err

}
