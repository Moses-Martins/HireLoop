package main

import (
	"io"
	"os"
	"fmt"
	"mime"
	"strings"
	"crypto/rand"
	"encoding/base64"
	"path/filepath"
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
)

type applyJob struct {
	ID uuid.UUID `json:"id"`
    ApplicantID uuid.UUID `json:"applicant_id"`
    JobID uuid.UUID `json:"job_id"`
    ResumeUrl string `json:"resume_url"`
    Status string `json:"status"`
}


func (cfg *apiConfig) applyForJobs(w http.ResponseWriter, req *http.Request) {
	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	ValidatedID, err := auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	respBodyInitial, err := cfg.DB.GetUserByID(req.Context(), ValidatedID)
	if err != nil {
    	w.WriteHeader(401)
		return
	}

	if respBodyInitial.Role != "applicant" {
		w.WriteHeader(400)
		w.Write([]byte("Only Applicants can apply for a Job"))
		return
	}

	vars := mux.Vars(req)         
    idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusNotFound)
		return
	}

	respBody, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		http.Error(w, "Cannot Retrieve Job", http.StatusNotFound)
		return
	}


	Applied, err := cfg.DB.ApplyJobs(req.Context(), database.ApplyJobsParams{
		ApplicantID: ValidatedID,
		JobID: respBody.ID,
		ResumeUrl: "Not added yet",
		Status: "Submitted",
	})
	if err != nil {
		http.Error(w, "Cannot Apply for Job", http.StatusNotFound)
		return
	}


	const maxMemory = 10 << 20
	req.ParseMultipartForm(maxMemory)
	
	file, header, err := req.FormFile("resume")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to parse form file"))
		return
	}

	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing Content-Type for resume"))
		return
	}


	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Content-Type"))
		return
	}

	if (mediatype != "application/pdf") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot upload content of that type"))
		return
	}

	parts := strings.Split(mediatype, "/")

	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)

	randomString := base64.RawURLEncoding.EncodeToString(randomBytes)

	filename := fmt.Sprintf("%s.%s", randomString, parts[1])
	path := filepath.Join(cfg.assetsRoot, "/", filename)


	file_dir, err := os.Create(path)
    if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot create file path"))
		return
    }
    defer file_dir.Close()

	_, err = io.Copy(file_dir, file)


	data_url := fmt.Sprintf("http://localhost:%s/assets/%s.%s", cfg.port, randomString, parts[1])
		

	err = cfg.DB.UpdateApplyJob(req.Context(), database.UpdateApplyJobParams{
		ResumeUrl: data_url,
		ID: Applied.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot update video"))
		return
	}

	updated, err := cfg.DB.GetApplyJobs(req.Context(), Applied.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot get resume"))
		return
	}

	Resp := applyJob{
		ID: updated.ID,
		ApplicantID: updated.ApplicantID,
		JobID: updated.JobID,
		ResumeUrl: updated.ResumeUrl,
		Status:  updated.Status,
	}

	data, err := json.Marshal(Resp)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
	
}