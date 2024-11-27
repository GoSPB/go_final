package models

const DateFormat = "20060102"
const SearchDateFormat = "02.01.2006"
const Limit = 50

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type IDResponse struct {
	ID string `json:"id"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}
