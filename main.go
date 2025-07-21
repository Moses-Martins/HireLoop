package main

import(
	"os"
	"log"
	"net/http"
	"database/sql"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/google/uuid"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv" 
	"github.com/Moses-Martins/HireLoop/internal/database"
)


type apiConfig struct {
	port string
	JwtSecret string
	DB *database.Queries
	GoogleOauthConfig *oauth2.Config
	assetsRoot       string
}

type resume struct {
	data      []byte
	mediaType string
}

var applicantResume = map[uuid.UUID]resume{}

 
func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("SECRET")
	assetsRoot := os.Getenv("ASSETS_ROOT")
	
	redirectURL :=  os.Getenv("GOOGLE_REDIRECT_URL")
	clientID :=     os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	scopes := []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		}
	endpoint := google.Endpoint



	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}


	dbQueries := database.New(db)

	apiCfg := apiConfig {
		port: port,
		DB: dbQueries,
		JwtSecret: jwtSecret,
		assetsRoot: assetsRoot,
		GoogleOauthConfig: &oauth2.Config{
			RedirectURL: redirectURL,
			ClientID:    clientID,
			ClientSecret: clientSecret,
			Scopes: scopes,
			Endpoint: endpoint,
		},
	}

	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("./assets"))
    router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))
	router.HandleFunc("/api/auth/register", apiCfg.CreateUsers).Methods("POST")
	router.HandleFunc("/api/auth/login", apiCfg.Login).Methods("POST")
	router.HandleFunc("/api/auth/google/login", apiCfg.GoogleLogin).Methods("GET")
	router.HandleFunc("/api/auth/google/callback", apiCfg.GoogleCallback).Methods("GET")
	router.HandleFunc("/api/auth/me", apiCfg.Me).Methods("GET")
	router.HandleFunc("/api/jobs", apiCfg.createJob).Methods("POST")
	router.HandleFunc("/api/jobs", apiCfg.getAllJobs).Methods("GET")
	router.HandleFunc("/api/jobs/{id}", apiCfg.getJobByID).Methods("GET")
	router.HandleFunc("/api/jobs/{id}", apiCfg.updateJobByID).Methods("PUT")
	router.HandleFunc("/api/jobs/{id}", apiCfg.deleteJobByID).Methods("DELETE")
	router.HandleFunc("/api/jobs/{id}/apply", apiCfg.applyForJobs).Methods("POST")
	router.HandleFunc("/api/employers/{id}/applications", apiCfg.getAllApp).Methods("GET")
	router.HandleFunc("/api/applications/{id}", apiCfg.viewSingleApp).Methods("GET")
	router.HandleFunc("/api/applications/{id}", apiCfg.deleteApplication).Methods("DELETE")
	

	srv := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}


	log.Printf("Serving on: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())


}