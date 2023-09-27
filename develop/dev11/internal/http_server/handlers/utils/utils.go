package utils

import (
	"calendar/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Response - response
type Response struct {
	Message string `json:"message,omitempty"`
}

// RenderJSON - respond with json and status code
func RenderJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// DTO - data transfer object
type DTO struct {
	UserId int    `json:"user_id"`
	Date   string `json:"date"`
}

// ValidateEvent - validate event
func ValidateEvent(event *model.Event) error {
	return validateEvent(event)
}

// ValidateDTO - validate DTO
func ValidateDTO(dto *DTO) error {
	return validateDTO(dto)
}

func validateEvent(event *model.Event) error {
	err := validateDTO(&DTO{UserId: event.UserId, Date: event.Date})
	if err != nil {
		return err
	}

	if event.Title == "" {
		return fmt.Errorf("invalid title: %s", event.Title)
	}

	return nil
}

func validateDTO(dto *DTO) error {

	if dto.UserId <= 0 {
		return fmt.Errorf("invalid user id: %d", dto.UserId)
	}
	_, err := time.Parse("2006-01-02", dto.Date)
	if err != nil {
		return fmt.Errorf("invalid date: %s", dto.Date)
	}

	return nil
}

// GetEvent - get event from request
func GetEvent(r *http.Request) (*model.Event, error) {
	var event *model.Event

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return nil, fmt.Errorf("invalid request")
	}
	return event, nil
}

// GetDTO - get DTO from request
func GetDTO(r *http.Request) (*DTO, error) {
	var dto *DTO

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("invalid request")
	}
	return dto, nil
}

// GetDTOFormParams - get DTO from url params
func GetDTOFormParams(r *http.Request) (*DTO, error) {
	dto := &DTO{}
	queryParams := r.URL.Query()
	id := queryParams.Get("user_id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %s", id)
	}
	dto.UserId = userId
	dto.Date = queryParams.Get("date")
	return dto, nil
}
