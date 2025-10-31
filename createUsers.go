package main

import (
	"encoding/json"
	"errors"
	"github.com/Moses-Martins/HireLoop/internal/auth"
	"github.com/Moses-Martins/HireLoop/internal/database"
	"net/http"
	"strings"
	"time"
	"regexp"
)

type AcceptEmail struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type UserShown struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
}


// CreateUsers godoc
// @Summary Register a new user
// @Description Registers a new user with email, password, name, and role (applicant or employer). Email is validated for proper format.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body AcceptEmail true "User registration info"
// @Success 201 {object} UserShown "User successfully registered"
// @Failure 400 {object} map[string]string "Invalid request, role, or email format"
// @Failure 500 {object} map[string]string "Internal server error or hashing failure"
// @Router /api/auth/register [post]
func (cfg *apiConfig) CreateUsers(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := AcceptEmail{}

	err := decoder.Decode(&params)
	if err != nil {
		Send(w, 500, nil, "Error decoding parameters")
		return
	}

	params.Email, err = validateEmail(params.Email)
	if err != nil {
		Send(w, 400, nil, err.Error())
		return
	}

	params.Password, err = auth.HashPassword(params.Password)
	if err != nil {
		Send(w, 500, nil, "Error hashing password")
		return
	}

	params.Role, err = validateRole(params.Role)
	if err != nil {
		Send(w, 400, nil, err.Error())
		return
	}

	userDb, err := cfg.DB.CreateUser(req.Context(), database.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.Password,
		Role:           params.Role,
	})
	if err != nil {
		Send(w, 500, nil, "Cannot create user")
		return
	}

	respBody := UserShown{
		CreatedAt: userDb.CreatedAt,
		Name:      userDb.Name,
		Email:     userDb.Email,
		Role:      userDb.Role,
	}

	Send(w, 201, respBody, "User registered")
}

func validateRole(role string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(role))

	switch normalized {
	case "employer", "applicant":
		return normalized, nil
	default:
		return "", errors.New("invalid role: must be 'employer' or 'applicant'")
	}
}



// validateEmail checks if the provided email string is valid
func validateEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	// Basic RFC 5322 email regex
	const emailRegexPattern = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegexPattern)

	if !re.MatchString(email) {
		return "", errors.New("invalid email format")
	}

	return email, nil
}
