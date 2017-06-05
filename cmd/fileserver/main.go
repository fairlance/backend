package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fairlance/backend/middleware"
	respond "gopkg.in/matryer/respond.v1"
)

var (
	folderPath = flag.String("folderPath", "/tmp/files", "Storage path.")
	port       = flag.Int("port", 3006, "Port.")
	secret     = flag.String("secret", "", "Secret.")
	opts       *respond.Options
)

func init() {
	flag.Parse()
	err := os.MkdirAll(*folderPath, 0755)
	if err != nil {
		log.Fatalf("error creating folder: %v", err)
	}

	// f, err := os.OpenFile("/var/log/fairlance/fileserver.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// log.SetOutput(f)
}

func main() {
	http.Handle("/file/", middleware.Chain(
		ensureMethod("GET"),
		middleware.JSONEnvelope,
		middleware.WithTokenFromHeader,
		middleware.AuthenticateTokenWithClaims(*secret),
	)(http.StripPrefix("/file/", http.FileServer(http.Dir(*folderPath)))))

	http.Handle("/upload", middleware.Chain(
		ensureMethod("POST"),
		middleware.JSONEnvelope,
		middleware.WithTokenFromHeader,
		middleware.AuthenticateTokenWithClaims(*secret),
	)(upload()))

	log.Printf("Listening on: %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func ensureMethod(method string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
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
}

func upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			respond.With(w, r, http.StatusMethodNotAllowed, errors.New("bad method, only POST is allowed"))
			return
		}

		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(err)
			respond.With(w, r, http.StatusBadRequest, err)
			return
		}
		defer file.Close()

		fileStartBuff := make([]byte, 512)
		read, err := file.Read(fileStartBuff)
		if err != nil {
			log.Println(err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		if read == 0 {
			err = errors.New("0 bytes were read for content type detecting")
			log.Println(err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		fileType := http.DetectContentType(fileStartBuff)

		file.Seek(0, 0)
		if err != nil {
			log.Println(err)
			respond.With(w, r, http.StatusInternalServerError, err)
			return
		}

		f, err := os.Create(*folderPath + "/" + header.Filename)
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

		respond.With(w, r, http.StatusOK, struct {
			Name string `json:"name"`
			URL  string `json:"url"`
			Type string `json:"type"`
		}{
			Name: header.Filename,
			URL:  "file/" + header.Filename,
			Type: fileType,
		})
	})
}

// func index() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Content-Type", "text/html")
// 		w.Write([]byte(`
//         <form enctype="multipart/form-data" action="/upload" method="post">
//             <input type="file" name="uploadfile" />
//             <input type="submit" value="upload" />
//         </form>
//         `))
// 	})
// }
