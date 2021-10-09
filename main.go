package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func init() {
	if _, err := os.Stat(Config.WorkingDir); os.IsNotExist(err) {
		err := os.Mkdir(Config.WorkingDir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	router := mux.NewRouter()

	http.HandleFunc("/", httpInterceptor(router))

	router.HandleFunc("/compile", compileEndpoint).Methods("POST")

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", Config.IP, Config.Port),
	}

	log.Fatal(server.ListenAndServe())

}

func httpInterceptor(router http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		auth := req.Header.Get("Auth")
		if auth == Config.AuthToken {
			router.ServeHTTP(rw, req)
		} else {
			rw.WriteHeader(http.StatusForbidden)
		}
	}
}
