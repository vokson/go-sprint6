package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

const htmlFileName = "index.html"

func HandleUpload(res http.ResponseWriter, req *http.Request) {
	file, header, err := req.FormFile("myFile")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	converted := []byte(service.ConvertText(string(data)))
	fileName := filepath.Join("tmp", time.Now().UTC().String()+filepath.Ext(header.Filename))

	if err := os.WriteFile(fileName, converted, 0755); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := res.Write(converted); err != nil {
		log.Printf("handle upload: write response: %v", err)
		return
	}

}

func HandleMain(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, htmlFileName)
}
