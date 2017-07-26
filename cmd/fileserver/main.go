package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"strings"

	"github.com/fairlance/backend/middleware"
	"github.com/google/uuid"
	respond "gopkg.in/matryer/respond.v1"
)

const _3MB = 3 * 1024 * 1024

func main() {
	log.SetFlags(log.Lshortfile)
	port := os.Getenv("PORT")
	secret := os.Getenv("SECRET")
	dataDir := os.Getenv("DATA_DIR")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("error creating folder: %v", err)
	}
	http.Handle("/public/file/", middleware.Chain(
		middleware.CORSHandler,
		middleware.JSONEnvelope,
		middleware.HTTPMethod("GET"),
	)(http.StripPrefix("/public/file/", http.FileServer(http.Dir(dataDir)))))
	http.Handle("/public/upload", middleware.Chain(
		middleware.CORSHandler,
		middleware.JSONEnvelope,
		middleware.WithTokenFromHeader,
		middleware.AuthenticateTokenWithUser(secret),
		middleware.HTTPMethod("POST"),
	)(upload(dataDir)))
	log.Printf("Listening on: %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func upload(dataDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, _3MB)
		file, header, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(fmt.Errorf("could not parse form: %v", err))
			respond.With(w, r, http.StatusBadRequest, fmt.Errorf("could not parse form: %v", err))
			return
		}
		defer file.Close()
		fileStartBuff := make([]byte, 512)
		read, err := file.Read(fileStartBuff)
		if err != nil {
			log.Println(fmt.Errorf("could not read bytes from file: %v", err))
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not read bytes from file: %v", err))
			return
		}
		if read == 0 {
			log.Println(fmt.Errorf("read 0 bytes from file: %v", err))
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("read 0 bytes from file: %v", err))
			return
		}
		fileType := http.DetectContentType(fileStartBuff)
		file.Seek(0, 0)
		uuid, err := uuid.NewUUID()
		ext := filepath.Ext(header.Filename)
		if ext == "" {
			ext = "." + strings.Split(fileType, "/")[1]
		}
		name := uuid.String() + ext
		if err != nil {
			log.Println(fmt.Errorf("could not generate uuid: %v", err))
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not generate uuid: %v", err))
			return
		}
		f, err := os.Create(dataDir + "/" + name)
		if err != nil {
			log.Println(fmt.Errorf("could not create destination file: %v", err))
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not create destination file: %v", err))
			return
		}
		defer f.Close()
		written, err := io.Copy(f, file)
		if err != nil {
			log.Println(fmt.Errorf("could not copy file to destination file: %v", err))
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("could not copy file to destination file: %v", err))
			return
		}
		if written == 0 {
			log.Println(fmt.Errorf("wrote 0 bytes to destination file: %v", err))
			respond.With(w, r, http.StatusInternalServerError, fmt.Errorf("wrote 0 bytes to destination file: %v", err))
			return
		}
		respond.With(w, r, http.StatusOK, struct {
			Name string `json:"name"`
			URL  string `json:"url"`
			Type string `json:"type"`
		}{
			Name: name,
			URL:  "file/" + name,
			Type: fileType,
		})
	})
}
