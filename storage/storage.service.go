package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mehmet.com/cors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const documentURLPath = "documents"

func handleDocuments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		documentList, err := getDocuments()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(documentList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("document")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		f, err := os.OpenFile(filepath.Join(DocumentDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()

		io.Copy(f, file)
		w.WriteHeader(201)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", documentURLPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fileName := urlPathSegments[1:][0]
	file, err := os.Open(filepath.Join(DocumentDirectory, fileName))
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fHeader := make([]byte, 512)
	file.Read(fHeader)
	fContentType := http.DetectContentType(fHeader)

	stat, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fSize := strconv.FormatInt(stat.Size(), 10)
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Length", fSize)
	file.Seek(0, 0)
	io.Copy(w, file)

}

func SetupRoutes(apiBasePath string) {
	storageHandler := http.HandlerFunc(handleDocuments)
	downloadHandler := http.HandlerFunc(handleDownload)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, documentURLPath), cors.Middleware(storageHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, documentURLPath), cors.Middleware(downloadHandler))

}
