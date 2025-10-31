package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"github.com/Moses-Martins/HireLoop/internal/database"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	_ "github.com/Moses-Martins/HireLoop/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title HireLoop API
// @version 1.0
// @description API for HireLoop job platform
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.hireloop.com/support
// @contact.email support@hireloop.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host hireloop.onrender.com
// @basePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
type apiConfig struct {
	port                string
	JwtSecret           string
	DB                  *database.Queries
	RegisterRedirectUrl string
	LoginRedirectUrl    string
	GoogleOauthConfig   *oauth2.Config
	assetsRoot          string
}

type resume struct {
	data      []byte
	mediaType string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("SECRET")
	assetsRoot := os.Getenv("ASSETS_ROOT")
	RegisterRedirectUrl := os.Getenv("REGISTER_REDIRECT_URL")
	LoginRedirectUrl := os.Getenv("LOGIN_REDIRECT_URL")
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
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

	apiCfg := apiConfig{
		port:                port,
		DB:                  dbQueries,
		JwtSecret:           jwtSecret,
		assetsRoot:          assetsRoot,
		RegisterRedirectUrl: RegisterRedirectUrl,
		LoginRedirectUrl:    LoginRedirectUrl,
		GoogleOauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			Endpoint:     endpoint,
		},
	}

	router := mux.NewRouter()

	// Serve static assets
	fs := http.FileServer(http.Dir("./assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Swagger UI endpoint
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Auth routes
	router.HandleFunc("/api/auth/register", apiCfg.CreateUsers).Methods("POST")
	router.HandleFunc("/api/auth/google/register", apiCfg.GoogleRegister).Methods("GET")
	router.HandleFunc("/api/auth/google/register/callback", apiCfg.RegisterCallback).Methods("GET")
	router.HandleFunc("/api/auth/login", apiCfg.Login).Methods("POST")
	router.HandleFunc("/api/auth/google/login", apiCfg.GoogleLogin).Methods("GET")
	router.HandleFunc("/api/auth/google/login/callback", apiCfg.LoginCallback).Methods("GET")
	router.HandleFunc("/api/auth/me", apiCfg.Me).Methods("GET")
	router.HandleFunc("/api/refresh", apiCfg.RefreshHandler).Methods("POST")
	router.HandleFunc("/api/revoke", apiCfg.RevokeHandler).Methods("POST")

	// Job routes
	router.HandleFunc("/api/jobs", apiCfg.createJob).Methods("POST")
	router.HandleFunc("/api/jobs", apiCfg.getAllJobs).Methods("GET")
	router.HandleFunc("/api/jobs/search", apiCfg.searchJobs).Methods("GET")
	router.HandleFunc("/api/jobs/filter", apiCfg.FilterJobs).Methods("GET")
	router.HandleFunc("/api/jobs/{id}", apiCfg.getJobByID).Methods("GET")
	router.HandleFunc("/api/jobs/{id}", apiCfg.updateJobByID).Methods("PUT")
	router.HandleFunc("/api/jobs/{id}", apiCfg.deleteJobByID).Methods("DELETE")
	router.HandleFunc("/api/jobs/{id}/apply", apiCfg.applyForJobs).Methods("POST")

	// Application routes
	router.HandleFunc("/api/employers/{id}/applications", apiCfg.getAllApp).Methods("GET")
	router.HandleFunc("/api/applications/{id}", apiCfg.viewSingleApp).Methods("GET")
	router.HandleFunc("/api/applications/{id}", apiCfg.deleteApplication).Methods("DELETE")

	// Start server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Serving on: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())
}
