package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
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

var ctx context.Context = context.Background()

func handle(writer http.ResponseWriter, req *http.Request, tmplQ *template.Template, tmplM *template.Template, rdb *redis.Client) {
	req.ParseForm()
	task := req.FormValue("c")
	numberS, _ := rdb.Get(ctx, task+":number").Result()
	number, _ := strconv.Atoi(numberS)
	sol := req.FormValue("solution")
	sol = strings.ToUpper(sol)
	fmt.Println([]byte(sol))
	if sol == "AHOJ" {
		mlok := gotoS{
			number,
			"50.8439058N, 14.2274044E",
			"Za márnicí."}
		tmplM.Execute(writer, mlok)
	} else {
		q := questionS{
			number,
			"BIPK",
		}
		tmplQ.Execute(writer, q)
	}
}

func main() {
	tmplQ := template.Must(template.ParseFiles("sifra.html"))
	tmplM := template.Must(template.ParseFiles("done.html"))

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		handle(writer, req, tmplQ, tmplM, rdb)
	})
	http.ListenAndServe("127.0.0.6:8080", nil)
}
