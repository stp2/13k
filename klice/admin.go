package main

import (
	"fmt"
	"net/http"
	"os"
)

func teams(writer http.ResponseWriter, req *http.Request) {
	iter := rdb.Scan(ctx, 0, "team/*/name", 0).Iterator()
	for iter.Next(ctx) {
		fmt.Fprintln(writer, iter.Val())
	}
}

func handleAdmin(writer http.ResponseWriter, req *http.Request) {
	user, pass, ok := req.BasicAuth()
	if ok {
		password, err := rdb.Get(ctx, "admin/user/"+user).Result()
		if err == nil && pass == password {
			path := req.URL.Path
			fmt.Println(path)
			if path[0] == '/' {
				path = path[1:]
			}
			fmt.Println(path)
			path = path[5:]
			if path != "" && path[0] == '/' {
				path = path[1:]
			}
			fmt.Println(path)
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
