package main

import (
	"net/http"
)

func handle(writer http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("makak")
	if err != nil {
		// set
		cookie := http.Cookie{
			Name:   "makak",
			Value:  "opice",
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(writer, &cookie)
		writer.Write([]byte("SET"))
	} else {
		// delete
		cookie.MaxAge = -1
		http.SetCookie(writer, cookie)
		writer.Write([]byte("UNSET"))
	}
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe("127.0.0.6:8080", nil)
}
