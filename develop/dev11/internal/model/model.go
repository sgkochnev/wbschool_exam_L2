package model

// Event - model
type Event struct {
	UserId  int    `json:"user_id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Details string `json:"details"`
}
