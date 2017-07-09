package importer

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/fairlance/backend/dispatcher"
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

type paymentHandler struct {
	db         *sql.DB
	dispatcher dispatcher.ApplicationDispatcher
}

func newPaymentHandler(db *sql.DB, applicationURL string) *paymentHandler {
	return &paymentHandler{
		db:         db,
		dispatcher: dispatcher.NewApplicationDispatcher(applicationURL),
	}
}

func (p *paymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var projects []project
	if r.Method == http.MethodPost {
		r.ParseForm()
		projectID, err := strconv.Atoi(r.Form.Get("project_id"))
		if err != nil {
			log.Printf("could not parse id: %v", err)
		}
		if err := p.dispatcher.SetProjectFunded(projectID); err != nil {
			log.Printf("could not set project status to funded: %v", err)
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
	rows, err := p.db.Query(selectSQL, step, offset)
	if err != nil {
		log.Printf("could not get projects: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var proj project
		err := rows.Scan(&proj.ID, &proj.Name, &proj.Status, &proj.ClientID)
		if err != nil {
			log.Printf("could not scan project: %v", err)
			return
		}
		projects = append(projects, proj)
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
