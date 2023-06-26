package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type questionS struct {
	Number   int
	Question template.HTML
}

func handle(writer http.ResponseWriter, req *http.Request, tmpl *template.Template) {
	req.ParseForm()
	sol := req.FormValue("solution")
	sol = strings.ToUpper(sol)
	fmt.Println([]byte(sol))
	if sol == "AHOJ" {
		writer.Write([]byte("Mlok"))
	} else {
		q := questionS{
			0,
			"BIPK",
		}
		tmpl.Execute(writer, q)
	}
}

func main() {
	tmpl := template.Must(template.ParseFiles("sifra.html"))

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		handle(writer, req, tmpl)
	})
	http.ListenAndServe("127.0.0.6:8080", nil)
}
