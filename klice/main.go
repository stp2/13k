package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
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

func handleSignIn(writer http.ResponseWriter, req *http.Request, rdb *redis.Client) {
	req.ParseForm()
	pass := req.FormValue("passphrase")
	_, err := rdb.Get(ctx, pass+":name").Result()
	if err == redis.Nil {
		body, _ := ioutil.ReadFile("signIn.html")
		fmt.Fprint(writer, string(body))
	} else {
		cookie := http.Cookie{
			Name:   "team",
			Value:  pass,
			Path:   "/",
			MaxAge: 36000,
		}
		task, err := req.Cookie("task")
		if err == nil {
			task.MaxAge = -1
			http.SetCookie(writer, task)
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, req, task.Value, 302)
	}
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
		_, err := req.Cookie("team")
		if err != nil {
			cookie := http.Cookie{
				Name:   "task",
				Value:  req.URL.Path,
				Path:   "/",
				MaxAge: 3600,
			}
			http.SetCookie(writer, &cookie)
			http.Redirect(writer, req, "/signin", 302)
		}
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
	http.HandleFunc("/signin", func(writer http.ResponseWriter, req *http.Request) {
		handleSignIn(writer, req, rdb)
	})
	http.ListenAndServe("127.0.0.6:8080", nil)
}
