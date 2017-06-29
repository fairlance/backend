package importer

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

const (
	step      = 10
	selectSQL = `SELECT id, name, status, client_id FROM projects ORDER BY id DESC LIMIT $1 OFFSET $2;`
)

type project struct {
	ID       string
	Name     string
	Status   string
	ClientID string
}

func paymentHandler(db *sql.DB, applicationURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var projects []project
		if r.Method == http.MethodPost {
			r.ParseForm()
			projectID, err := strconv.Atoi(r.Form.Get("project_id"))
			if err != nil {
				log.Printf("could not parse id: %v", err)
			}
			url := fmt.Sprintf("%s/private/project/%d/fund", applicationURL, projectID)
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Printf("could not create request: %v", err)
				return
			}
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				log.Printf("could not do request: %v", err)
				return
			}
			if response.StatusCode != http.StatusOK {
				log.Printf("bad status code when changing project status: %d", response.StatusCode)
				return
			}
			http.Redirect(w, r, "/api/importer/payment?"+r.Form.Get("get_params"), http.StatusSeeOther)
			return
		}
		var offset int
		if r.URL.Query().Get("offset") != "" {
			o, err := strconv.Atoi(r.URL.Query().Get("offset"))
			if err != nil {
				log.Printf("could not parse offset: %v", err)
				offset = 0
			}
			offset = o
		}
		rows, err := db.Query(selectSQL, step, offset)
		if err != nil {
			log.Printf("could not get projects: %v", err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var p project
			err := rows.Scan(&p.ID, &p.Name, &p.Status, &p.ClientID)
			if err != nil {
				log.Printf("could not scan project: %v", err)
				return
			}
			projects = append(projects, p)
		}

		main := MustAsset("templates/payment.html")
		tmpl, err := template.New("payment").Funcs(funcMap).Parse(string(main))
		if err != nil {
			log.Fatal(err)
		}
		if err := tmpl.Execute(w, struct {
			Projects  []project
			Offset    int
			Step      int
			GETParams string
		}{
			projects,
			offset,
			step,
			r.URL.RawQuery,
		}); err != nil {
			log.Fatal(err)
		}
		return
	})
}

var funcMap = template.FuncMap{
	"prev": func(a, b int) int {
		prev := a - b
		if prev >= 0 {
			return prev
		}
		return 0
	},
	"next": func(a, b int) int {
		return a + b
	},
}
