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
}

var AtmplT *template.Template

func getTeamName(path string) string {
	path, _ = strings.CutPrefix(path, "team/")
	path, _ = strings.CutSuffix(path, "/name")
	return path
}

func teams(writer http.ResponseWriter, req *http.Request) {
	var teams []AteamsT = make([]AteamsT, 0)
	var teamPass []string = make([]string, 0)

	iter := rdb.Scan(ctx, 0, "team/*/name", 0).Iterator()
	for iter.Next(ctx) {
		teamPass = append(teamPass, getTeamName(iter.Val()))
	}
	for _, t := range teamPass {
		name, _ := rdb.Get(ctx, "team/"+t+"/name").Result()
		tier, _ := rdb.Get(ctx, "team/"+t+"/tier").Result()
		lastS, _ := rdb.Get(ctx, "team/"+t+"/last").Result()
		last, _ := strconv.Atoi(lastS)
		teams = append(teams, AteamsT{name, t, tier, last})
	}
	AtmplT.Execute(writer, teams)
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
			path = path[5:]
			if path != "" && path[0] == '/' {
				path = path[1:]
			}
			switch path {
			case "teams":
				teams(writer, req)
			default:
				body, _ := os.ReadFile("admin.html")
				fmt.Fprint(writer, string(body))
			}
			return
		}
	}
	writer.Header().Add("WWW-Authenticate", `Basic realm="Přihlaš se!"`)
	writer.WriteHeader(http.StatusUnauthorized)
	return
}
