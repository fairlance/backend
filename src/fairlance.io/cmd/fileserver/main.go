package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

var folderPath = "/tmp/files"

func init() {
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		log.Fatalf("error creating folder: %v", err)
	}
}

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(`
        <form enctype="multipart/form-data" action="/upload" method="post">
            <input type="file" name="uploadfile" />
            <input type="submit" value="upload" />
        </form>
        `))
	}))
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(folderPath))))
	http.Handle("/upload", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" { //http.MethodPost {
			w.Write([]byte("bad method, only POST is allowed"))
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		f, err := os.Create(folderPath + "/" + handler.Filename)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer f.Close()

		written, err := io.Copy(f, file)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if written == 0 {
			err = errors.New("0 bytes were written")
			log.Println(err)
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte("file " + handler.Filename + " sucessfully uploaded\n"))
		w.Write([]byte("you can see it at /file/" + handler.Filename + ""))
	}))
	http.ListenAndServe(":3006", nil)
}
