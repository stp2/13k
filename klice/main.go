package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
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

type teamS struct {
	Name string
	Tier string
	Last int
}

var ctx context.Context = context.Background()
var rdb *redis.Client

func getTask(req *http.Request) (task string) {
	qr := req.URL.Path
	if qr[0] == '/' {
		qr = qr[1:]
	}
	team, _ := req.Cookie("team")
	tier, _ := rdb.Get(ctx, "team/"+team.Value+"/tier").Result()
	task, _ = rdb.Get(ctx, qr+"/tier/"+tier).Result()
	return
}

func handleCipher(writer http.ResponseWriter, req *http.Request, tmpl *template.Template) {
	task := getTask(req)
	numberS, _ := rdb.Get(ctx, task+"/number").Result()
	number, _ := strconv.Atoi(numberS)
	cipher, _ := rdb.Get(ctx, task+"/cipher").Result()
	q := questionS{
		number,
		template.HTML(cipher),
	}
	tmpl.Execute(writer, q)
}

func handleMlok(writer http.ResponseWriter, req *http.Request, tmpl *template.Template) {
	task := getTask(req)
	numberS, _ := rdb.Get(ctx, task+"/number").Result()
	number, _ := strconv.Atoi(numberS)
	// incr last
	team, _ := req.Cookie("team")
	rdb.Set(ctx, "team/"+team.Value+"/last", number, 0)
	// get next
	next, _ := rdb.Get(ctx, task+"/next").Result()
	position, _ := rdb.Get(ctx, next+"/position").Result()
	help, _ := rdb.Get(ctx, next+"/help").Result()
	mlok := gotoS{
		number,
		template.HTML(position),
		template.HTML(help)}
	tmpl.Execute(writer, mlok)
}

func handleSignIn(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	pass := req.FormValue("passphrase")
	_, err := rdb.Get(ctx, "team/"+pass+"/name").Result()
	if err == redis.Nil {
		body, _ := os.ReadFile("signIn.html")
		fmt.Fprint(writer, string(body))
	} else {
		cookie := http.Cookie{
			Name:   "team",
			Value:  pass,
			Path:   "/",
			MaxAge: 36000,
		}
		qr, err := req.Cookie("qr")
		if err == nil {
			qr.MaxAge = -1
			http.SetCookie(writer, qr)
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, req, qr.Value, http.StatusFound)
	}
}

func solved(req *http.Request) bool {
	// last solved task
	team, _ := req.Cookie("team")
	lastS, _ := rdb.Get(ctx, "team/"+team.Value+"/last").Result()
	last, _ := strconv.Atoi(lastS)
	// task number
	task := getTask(req)
	numberS, _ := rdb.Get(ctx, task+"/number").Result()
	number, _ := strconv.Atoi(numberS)
	if last >= number {
		return true
	} else {
		return false
	}
}

func isSignIn(writer http.ResponseWriter, req *http.Request) bool {
	_, err := req.Cookie("team")
	if err != nil { // need login
		cookie := http.Cookie{ // qr url
			Name:   "qr",
			Value:  req.URL.Path,
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, req, "/signin", http.StatusFound)
		return false
	} else {
		return true
	}
}

func handleTeam(writer http.ResponseWriter, req *http.Request, template template.Template) {
	team, _ := req.Cookie("team")
	name, _ := rdb.Get(ctx, "team/"+team.Value+"/name").Result()
	tier, _ := rdb.Get(ctx, "team/"+team.Value+"/tier").Result()
	lastS, _ := rdb.Get(ctx, "team/"+team.Value+"/last").Result()
	last, _ := strconv.Atoi(lastS)
	info := teamS{
		name,
		tier,
		last,
	}
	template.Execute(writer, info)
}

func main() {
	tmplQ := template.Must(template.ParseFiles("cipher.html"))
	tmplM := template.Must(template.ParseFiles("done.html"))
	tmplT := template.Must(template.ParseFiles("team.html"))

	// Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	http.HandleFunc("/qr/", func(writer http.ResponseWriter, req *http.Request) {
		if isSignIn(writer, req) {
			req.ParseForm()
			task := getTask(req)
			solOk, _ := rdb.Get(ctx, task+"/solution").Result()
			sol := req.FormValue("solution")
			sol = strings.ToUpper(sol)
			if sol == solOk || solved(req) {
				handleMlok(writer, req, tmplM)
			} else {
				handleCipher(writer, req, tmplQ)
			}
		}
	})
	http.HandleFunc("/signin", func(writer http.ResponseWriter, req *http.Request) {
		handleSignIn(writer, req)
	})
	http.HandleFunc("/team", func(writer http.ResponseWriter, req *http.Request) {
		if isSignIn(writer, req) {
			handleTeam(writer, req, *tmplT)
		}
	})
	http.ListenAndServe("127.0.0.6:8080", nil)
}
