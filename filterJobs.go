package main

import (
	"database/sql"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"net/http"
	"strconv"
)

// FilterJobs godoc
// @Summary Filter jobs
// @Description Filter job listings by location, type, and salary range. Requires authentication.
// @Tags jobs
// @Produce json
// @Security ApiKeyAuth
// @Param location query string false "Filter by location"
// @Param type query string false "Filter by job type (full-time, part-time, freelance, contract)"
// @Param salary_min query number false "Minimum salary"
// @Param salary_max query number false "Maximum salary"
// @Success 200 {array} jobs "Filtered list of jobs"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 401 {object} map[string]string "Invalid or missing token"
// @Failure 404 {object} map[string]string "Cannot retrieve filtered jobs"
// @Router /api/jobs/filter [get]
func (cfg *apiConfig) FilterJobs(w http.ResponseWriter, req *http.Request) {
	location := req.URL.Query().Get("location")
	jobType := req.URL.Query().Get("type")
	salaryMinStr := req.URL.Query().Get("salary_min")
	salaryMaxStr := req.URL.Query().Get("salary_max")

	var salaryMin float32 = 1
	var salaryMax float32 = 1
	var err error

	if salaryMinStr != "" {
		conv, err := strconv.ParseFloat(salaryMinStr, 32)
		if err != nil {
			Send(w, 400, nil, "Invalid salary_min")
			return
		}
		salaryMin = float32(conv)
	}

	if salaryMaxStr != "" {
		conv, err := strconv.ParseFloat(salaryMaxStr, 32)
		if err != nil {
			Send(w, 400, nil, "Invalid salary_max")
			return
		}
		salaryMax = float32(conv)
	}

	token_string, err := auth.GetBearerToken(req.Header)
	if err != nil {
		Send(w, 401, nil, "Invalid or missing token")
		return
	}

	_, err = auth.ValidateJWT(token_string, cfg.JwtSecret)
	if err != nil {
		Send(w, 401, nil, "Invalid or missing token")
		return
	}

	JobDb, err := cfg.DB.FiltersJobs(req.Context(), database.FiltersJobsParams{
		Column1: sql.NullString{
			String: location,
			Valid:  true,
		},
		Column2: jobType,
		Column3: salaryMin,
		Column4: salaryMax,
	})
	if err != nil {
		Send(w, 404, nil, "Cannot retrieve filtered jobs")
		return
	}

	respBody := make([]jobs, 0, len(JobDb))
	for _, dbjob := range JobDb {
		jobResp := jobs{
			ID:          dbjob.ID,
			CreatedAt:   dbjob.CreatedAt,
			UpdatedAt:   dbjob.UpdatedAt,
			Title:       dbjob.Title,
			Description: dbjob.Description,
			Location:    dbjob.Location,
			Type:        dbjob.Type,
			Salary:      dbjob.Salary,
			EmployerID:  dbjob.EmployerID,
		}
		respBody = append(respBody, jobResp)
	}

	Send(w, 200, respBody, "Jobs filtered")
}
