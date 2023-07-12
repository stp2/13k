package main

import (
	"context"
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

func getTask(req *http.Request) (task string) {
	task = req.URL.Path
	if task[0] == '/' {
		task = task[1:]
	}
	return
}

func handleCipher(writer http.ResponseWriter, req *http.Request, tmpl *template.Template, rdb *redis.Client) {
	task := getTask(req)
	numberS, _ := rdb.Get(ctx, task+":number").Result()
	number, _ := strconv.Atoi(numberS)
	cipher, _ := rdb.Get(ctx, task+":cipher").Result()
	q := questionS{
		number,
		template.HTML(cipher),
	}
	tmpl.Execute(writer, q)
}

func handleMlok(writer http.ResponseWriter, req *http.Request, tmpl *template.Template, rdb *redis.Client) {
	task := getTask(req)
	numberS, _ := rdb.Get(ctx, task+":number").Result()
	number, _ := strconv.Atoi(numberS)
	position, _ := rdb.Get(ctx, task+":position").Result()
	help, _ := rdb.Get(ctx, task+":help").Result()
	mlok := gotoS{
		number,
		template.HTML(position),
		template.HTML(help)}
	tmpl.Execute(writer, mlok)
}

func main() {
	tmplQ := template.Must(template.ParseFiles("cipher.html"))
	tmplM := template.Must(template.ParseFiles("done.html"))

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		task := getTask(req)
		solOk, _ := rdb.Get(ctx, task+":solution").Result()
		sol := req.FormValue("solution")
		sol = strings.ToUpper(sol)
		if sol == solOk {
			handleMlok(writer, req, tmplM, rdb)
		} else {
			handleCipher(writer, req, tmplQ, rdb)
		}
	})
	http.ListenAndServe("127.0.0.6:8080", nil)
}
