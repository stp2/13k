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

func getTask(req *http.Request, rdb *redis.Client) (task string) {
	task = req.URL.Path
	if task[0] == '/' {
		task = task[1:]
	}
	team, _ := req.Cookie("team")
	tier, _ := rdb.Get(ctx, team.Value+":tier").Result()
	task, _ = rdb.Get(ctx, task+":"+tier).Result()
	return
}

func handleCipher(writer http.ResponseWriter, req *http.Request, tmpl *template.Template, rdb *redis.Client) {
	task := getTask(req, rdb)
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
	task := getTask(req, rdb)
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

func isValidTeam(team string, rdb *redis.Client) bool {
	_, err := rdb.Get(ctx, team+":name").Result()
	if err == redis.Nil {
		return false
	}
	return true
}

func isSignIn(writer http.ResponseWriter, req *http.Request, rdb *redis.Client) bool {
	// send form
	pass := req.FormValue("passphrase")
	if isValidTeam(pass, rdb) {
		cookie := http.Cookie{
			Name:   "team",
			Value:  pass,
			Path:   "/",
			MaxAge: 36000,
		}
		http.SetCookie(writer, &cookie)
		req.AddCookie(&cookie)
		return true
	}
	// cookie
	team, err := req.Cookie("team")
	if err != nil {
		return false
	}
	if isValidTeam(team.Value, rdb) {
		return true
	}
	return false
}

func handleSignIn(writer http.ResponseWriter) {
	body, _ := ioutil.ReadFile("signIn.html")
	fmt.Fprint(writer, string(body))
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
		if isSignIn(writer, req, rdb) {
			task := getTask(req, rdb)
			solOk, _ := rdb.Get(ctx, task+":solution").Result()
			sol := req.FormValue("solution")
			sol = strings.ToUpper(sol)
			if sol == solOk {
				handleMlok(writer, req, tmplM, rdb)
			} else {
				handleCipher(writer, req, tmplQ, rdb)
			}
		} else {
			handleSignIn(writer)
		}
	})
	http.ListenAndServe("127.0.0.6:8080", nil)
}
