package main

import(
	"log"
	"net/http"
	"github.com/gorilla/mux"
)


func main() {

	router := mux.NewRouter()
	//router.HandleFunc("/api/health", )

	srv := &http.Server{
		Handler: router,
		Addr: ":8080",
	}

	log.Fatal(srv.ListenAndServe())


}