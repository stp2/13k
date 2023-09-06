package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type AteamsT struct {
	Name string
	Pass string
	Tier string
	Last int
	Next string
}

type ATaskT struct {
	Number   string
	Cipher   template.HTML
	Solution string
	Position string
	Help     string
	Link     template.HTML
}

type APathT struct {
	Tier string
	Path []ATaskT
}

var AtmplT *template.Template
var AtmplP *template.Template

func getTeamName(path string) string {
	path, _ = strings.CutPrefix(path, "team/")
	path, _ = strings.CutSuffix(path, "/name")
	return path
}

func getTeams() (teams []string) {
	iter := rdb.Scan(ctx, 0, "team/*/name", 0).Iterator()
	for iter.Next(ctx) {
		teams = append(teams, getTeamName(iter.Val()))
	}
	return
}

func resetLast(writer http.ResponseWriter, req *http.Request) {
	teams := getTeams()
	for _, t := range teams {
		rdb.Set(ctx, "team/"+t+"/last", 0, -1)
	}
	http.Redirect(writer, req, "/admin/teams", http.StatusFound)
}

func teams(writer http.ResponseWriter, req *http.Request) {
	var teams []AteamsT = make([]AteamsT, 0)

	teamPass := getTeams()
	// get data
	for _, t := range teamPass {
		name, _ := rdb.Get(ctx, "team/"+t+"/name").Result()
		tier, _ := rdb.Get(ctx, "team/"+t+"/tier").Result()
		lastS, _ := rdb.Get(ctx, "team/"+t+"/last").Result()
		last, _ := strconv.Atoi(lastS)
		teams = append(teams, AteamsT{name, t, tier, last, ""})
	}
	AtmplT.Execute(writer, teams)
}

func getTierPath(tier string) (path []string) {
	qr, _ := rdb.Get(ctx, "start/"+tier).Result()
	for qr[:3] != "end" {
		path = append(path, qr)
		task, _ := rdb.Get(ctx, qr+"/tier/"+tier).Result()
		qr, _ = rdb.Get(ctx, task+"/next").Result()
	}
	path = append(path, qr)
	return
}

func handlePath(writer http.ResponseWriter, tier string) {
	var pathS APathT
	pathS.Tier = tier
	path := getTierPath(tier)
	for _, qr := range path {
		task, _ := rdb.Get(ctx, qr+"/tier/"+tier).Result()
		number, _ := rdb.Get(ctx, task+"/number").Result()
		cipher, _ := rdb.Get(ctx, task+"/cipher").Result()
		solution, _ := rdb.Get(ctx, task+"/solution").Result()
		position, _ := rdb.Get(ctx, qr+"/position").Result()
		help, _ := rdb.Get(ctx, qr+"/help").Result()
		link := template.HTML("https://klice.h21.fun/" + qr)
		if qr[:3] == "end" {
			number = ""
			link = ""
		}
		pathS.Path = append(pathS.Path, ATaskT{number, template.HTML(cipher), solution, position, help, link})
	}
	AtmplP.Execute(writer, pathS)
}

func handleAdmin(writer http.ResponseWriter, req *http.Request) {
	user, pass, ok := req.BasicAuth()
	if ok {
		password, err := rdb.Get(ctx, "admin/user/"+user).Result()
		if err == nil && pass == password {
			path := req.URL.Path
			if path[0] == '/' {
				path = path[1:]
			}
			path, _ = strings.CutPrefix(path, "admin")
			if path != "" && path[0] == '/' {
				path = path[1:]
			}
			switch {
			case path == "teams":
				teams(writer, req)
			case strings.Split(path, "/")[0] == "tasks":
				tier, _ := strings.CutPrefix(path, "tasks/")
				handlePath(writer, tier)
			default:
				body, _ := os.ReadFile("admin.html")
				fmt.Fprint(writer, string(body))
			}
			return
		}
	}
	writer.Header().Add("WWW-Authenticate", `Basic realm="Přihlaš se!"`)
	writer.WriteHeader(http.StatusUnauthorized)
}
