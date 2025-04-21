package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/dapoadedire/fem_project/internal/store"
	"github.com/dapoadedire/fem_project/internal/utils"
)

type registeredUserRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Bio            string `json:"bio"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(reg *registeredUserRequest) error {
	const (
		minUsernameLength = 3
		maxUsernameLength = 30
		minPasswordLength = 8
		minNameLength     = 2
	)

	// Trim input to avoid leading/trailing whitespace issues
	reg.Username = strings.TrimSpace(reg.Username)
	reg.Email = strings.TrimSpace(reg.Email)
	reg.FirstName = strings.TrimSpace(reg.FirstName)
	reg.LastName = strings.TrimSpace(reg.LastName)

	// Username validation
	if reg.Username == "" {
		return errors.New("username is required")
	}
	if len(reg.Username) < minUsernameLength {
		return errors.New(fmt.Sprintf("username must be at least %d characters long", minUsernameLength))
	}
	if len(reg.Username) > maxUsernameLength {
		return errors.New(fmt.Sprintf("username must be at most %d characters long", maxUsernameLength))
	}

	// Email validation
	if reg.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(reg.Email); err != nil {
		return errors.New("invalid email address")
	}
	

	// Password validation
	if reg.Password == "" {
		return errors.New("password is required")
	}
	if len(reg.Password) < minPasswordLength {
		return errors.New(fmt.Sprintf("password must be at least %d characters long", minPasswordLength))
	}

	// First name validation
	if len(reg.FirstName) < minNameLength {
		return errors.New(fmt.Sprintf("first name must be at least %d characters long", minNameLength))
	}

	// Last name validation
	if len(reg.LastName) < minNameLength {
		return errors.New(fmt.Sprintf("last name must be at least %d characters long", minNameLength))
	}

	return nil
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var regRequest registeredUserRequest
	err := json.NewDecoder(r.Body).Decode(&regRequest)
	if err != nil {
		h.logger.Printf("ERROR: decoding register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	err = h.validateRegisterRequest(&regRequest)
	if err != nil {
		h.logger.Printf("ERROR: validating register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username:       regRequest.Username,
		Email:          regRequest.Email,
		Bio:            regRequest.Bio,
		FirstName:      regRequest.FirstName,
		LastName:       regRequest.LastName,
		ProfilePicture: regRequest.ProfilePicture,
	}
	err = user.PasswordHash.SetPassword(regRequest.Password)
	if err != nil {
		h.logger.Printf("ERROR: setting password hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "user created successfully"})
}
