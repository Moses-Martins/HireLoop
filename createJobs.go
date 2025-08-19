package main

import (
	"encoding/json"
	"errors"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type job struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	Type        string  `json:"type"`
	Salary      float32 `json:"salary"`
}

type jobs struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Type        string    `json:"type"`
	Salary      float32   `json:"salary"`
	EmployerID  uuid.UUID `json:"employer_id"`
}

func (cfg *apiConfig) createJob(w http.ResponseWriter, req *http.Request) {

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

	if respBodyInitial.Role != "employer" {
		Send(w, 400, nil, "Only Employers can Create a Job")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := job{}

	err = decoder.Decode(&params)
	if err != nil {
		Send(w, 500, nil, "Error decoding parameters")
		return
	}

	params.Salary, err = validateFloat(params.Salary)
	if err != nil {
		Send(w, 400, nil, err.Error())
		return
	}

	params.Type, err = validateJobType(params.Type)
	if err != nil {
		Send(w, 400, nil, err.Error())
		return
	}

	jobDb, err := cfg.DB.CreateJobs(req.Context(), database.CreateJobsParams{
		Title:       params.Title,
		Description: params.Description,
		Location:    params.Location,
		Type:        params.Type,
		Salary:      params.Salary,
		EmployerID:  ValidatedID,
	})
	if err != nil {
		Send(w, 404, nil, "Cannot Create Job")
		return
	}

	respBody := jobs{
		ID:          jobDb.ID,
		CreatedAt:   jobDb.CreatedAt,
		UpdatedAt:   jobDb.UpdatedAt,
		Title:       jobDb.Title,
		Description: jobDb.Description,
		Location:    jobDb.Location,
		Type:        jobDb.Type,
		Salary:      jobDb.Salary,
		EmployerID:  jobDb.EmployerID,
	}

	Send(w, 201, respBody, "Job created")

}

func validateFloat(value interface{}) (float32, error) {
	switch v := value.(type) {
	case float32:
		return v, nil
	case float64:
		return float32(v), nil
	case int, int8, int16, int32, int64:
		return float32(reflect.ValueOf(v).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float32(reflect.ValueOf(v).Uint()), nil
	default:
		return 0, errors.New("invalid value: must be a number")
	}
}

func validateJobType(jobType string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(jobType))

	switch normalized {
	case "full-time", "part-time", "freelance", "contract":
		return normalized, nil
	default:
		return "", errors.New("invalid job type: must be 'full-time', 'part-time', 'freelance', or 'contract'")
	}
}
