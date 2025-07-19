package main

import(
	"os"
	"log"
	"net/http"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv" 
	"github.com/Moses-Martins/HireLoop/internal/database"
)


type apiConfig struct {

	DB *database.Queries
	
}


func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("PORT")




	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}


	dbQueries := database.New(db)

	apiCfg := apiConfig {
		DB: dbQueries,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/users", apiCfg.CreateUsers).Methods("POST")



	srv := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}


	log.Printf("Serving on: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())


}