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

type gotoS struct {
	Number int
	Goto   template.HTML
	Help   template.HTML
}

func handle(writer http.ResponseWriter, req *http.Request, tmplQ *template.Template, tmplM *template.Template) {
	req.ParseForm()
	sol := req.FormValue("solution")
	sol = strings.ToUpper(sol)
	fmt.Println([]byte(sol))
	if sol == "AHOJ" {
		mlok := gotoS{
			0,
			"50.8439058N, 14.2274044E",
			"Za márnicí."}
		tmplM.Execute(writer, mlok)
	} else {
		q := questionS{
			0,
			"BIPK",
		}
		tmplQ.Execute(writer, q)
	}
}

func main() {
	tmplQ := template.Must(template.ParseFiles("sifra.html"))
	tmplM := template.Must(template.ParseFiles("done.html"))

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		handle(writer, req, tmplQ, tmplM)
	})
	http.ListenAndServe("127.0.0.6:8080", nil)
}
