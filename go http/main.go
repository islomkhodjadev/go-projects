package main

import (
	"fmt"
	"net/http"
)

type app struct {
	addr string
}

func (s *app) getUserDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is handler for users "))
}
func (s *app) postNewUserCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is handler for users post creation "))
}
func (s *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/":
			fmt.Fprintf(w, "Hello,")
		case "/about":
			fmt.Fprintf(w, "this is about the product")
		}
	}

}

func main() {
	s := &app{addr: ":8080"}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}
	mux.HandleFunc("POST /users", s.postNewUserCreate)

	mux.HandleFunc("GET /users", s.getUserDataHandler)
	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
