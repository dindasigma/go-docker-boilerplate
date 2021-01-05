package application

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	router = mux.NewRouter()
)

func Run(addr string) {
	initializeRoutes()
	fmt.Printf("Listening to port %s", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}
