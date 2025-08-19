package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type applyJob struct {
	ID          uuid.UUID `json:"id"`
	ApplicantID uuid.UUID `json:"applicant_id"`
	JobID       uuid.UUID `json:"job_id"`
	ResumeUrl   string    `json:"resume_url"`
	Status      string    `json:"status"`
}

func (cfg *apiConfig) applyForJobs(w http.ResponseWriter, req *http.Request) {
	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		Send(w, 401, nil, "Invalid or missing token")
		return
	}

	ValidatedID, err := auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		Send(w, 401, nil, "Invalid or missing token")
		return
	}

	respBodyInitial, err := cfg.DB.GetUserByID(req.Context(), ValidatedID)
	if err != nil {
		Send(w, 400, nil, "Cannot be processed")
		return
	}

	if respBodyInitial.Role != "applicant" {
		Send(w, 400, nil, "Only Applicants can apply for a Job")
		return
	}

	vars := mux.Vars(req)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		Send(w, 404, nil, "Invalid UUID")
		return
	}

	respBody, err := cfg.DB.GetJobsByID(req.Context(), id)
	if err != nil {
		Send(w, 404, nil, "Cannot Retrieve Job")
		return
	}

	Applied, err := cfg.DB.ApplyJobs(req.Context(), database.ApplyJobsParams{
		ApplicantID: ValidatedID,
		JobID:       respBody.ID,
		ResumeUrl:   "Not added yet",
		Status:      "Submitted",
	})
	if err != nil {
		Send(w, 404, nil, "Cannot Apply for Job")
		return
	}

	const maxMemory = 10 << 20
	req.ParseMultipartForm(maxMemory)

	file, header, err := req.FormFile("resume")
	if err != nil {
		Send(w, 400, nil, "Unable to parse form file")
		return
	}

	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		Send(w, 400, nil, "Missing Content-Type for resume")
		return
	}

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		Send(w, 400, nil, "Invalid Content-Type")
		return
	}

	if mediatype != "application/pdf" {
		Send(w, 400, nil, "Cannot upload content of that type")
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
		Send(w, 500, nil, "Cannot create file path")
		return
	}
	defer file_dir.Close()

	_, err = io.Copy(file_dir, file)

	data_url := fmt.Sprintf("http://localhost:%s/assets/%s.%s", cfg.port, randomString, parts[1])

	err = cfg.DB.UpdateApplyJob(req.Context(), database.UpdateApplyJobParams{
		ResumeUrl: data_url,
		ID:        Applied.ID,
	})
	if err != nil {
		Send(w, 500, nil, "Cannot update video")
		return
	}

	updated, err := cfg.DB.GetApplyJobs(req.Context(), Applied.ID)
	if err != nil {
		Send(w, 500, nil, "Cannot get resume")
		return
	}

	Resp := applyJob{
		ID:          updated.ID,
		ApplicantID: updated.ApplicantID,
		JobID:       updated.JobID,
		ResumeUrl:   updated.ResumeUrl,
		Status:      updated.Status,
	}

	Send(w, 201, Resp, "Job application submitted")

}
