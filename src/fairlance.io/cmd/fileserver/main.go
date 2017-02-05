package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	respond "gopkg.in/matryer/respond.v1"

	app "fairlance.io/application"
)

var (
	folderPath = flag.String("folderPath", "/tmp/files", "Storage path.")
	port       = flag.String("port", "", "Port.")
	secret     = flag.String("secret", "", "Secret.")
	opts       *respond.Options
)

func init() {
	flag.Parse()
	err := os.MkdirAll(*folderPath, 0755)
	if err != nil {
		log.Fatalf("error creating folder: %v", err)
	}

	f, err := os.OpenFile("/var/log/fairlance/fileserver.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)

	opts = &respond.Options{
		Before: func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
			dataEnvelope := map[string]interface{}{"code": status}
			if err, ok := data.(error); ok {
				dataEnvelope["error"] = err.Error()
				dataEnvelope["success"] = false
			} else {
				dataEnvelope["data"] = data
				dataEnvelope["success"] = true
			}
			return status, dataEnvelope
		},
	}
}

func main() {
	http.Handle("/",
		opts.Handler(ensureMethod("GET", index())))

	http.Handle("/file/",
		opts.Handler(ensureMethod("GET", app.WithTokenFromHeader(
			app.AuthenticateTokenWithClaims(
				*secret, http.StripPrefix("/file/", http.FileServer(http.Dir(*folderPath))))))))

	http.Handle("/upload",
		opts.Handler(ensureMethod("POST", app.WithTokenFromHeader(
			app.AuthenticateTokenWithClaims(*secret, upload())))))

	http.ListenAndServe(":"+*port, nil)
}

func ensureMethod(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", r.Method+",OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		if r.Method == "OPTIONS" {
			// Stop here for a Preflighted OPTIONS request.
			return
		} else if r.Method != method {
			respond.With(w, r, http.StatusMethodNotAllowed, fmt.Errorf("bad method, only %s is allowed", method))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			respond.With(w, r, http.StatusMethodNotAllowed, errors.New("bad method, only POST is allowed"))
			return
		}

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(err)
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		defer file.Close()

		f, err := os.Create(*folderPath + "/" + handler.Filename)
		if err != nil {
			log.Println(err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}
		defer f.Close()

		written, err := io.Copy(f, file)
		if err != nil {
			log.Println(err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		if written == 0 {
			err = errors.New("0 bytes were written")
			log.Println(err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		respond.With(w, r, http.StatusOK, "file/"+handler.Filename+" sucessfully uploaded")
	})
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(`
        <form enctype="multipart/form-data" action="/upload" method="post">
            <input type="file" name="uploadfile" />
            <input type="submit" value="upload" />
        </form>
        `))
	})
}
