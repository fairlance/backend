package importer

import (
	"html/template"
	"log"
	"net/http"
)

type MessagesHandler struct{}

func NewMessagesHandler() *MessagesHandler {
	return &MessagesHandler{}
}

func (handler *MessagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	main := MustAsset("templates/messages.html")

	tmpl, err := template.New("messages").Parse(string(main))
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
	return
}
