package model

type OpenLibrarySubjectsResponse struct {
	// Key         string `json:"key"`
	Name        string `json:"name"`
	// SubjectType string `json:"subject_type"`
	WorkCount   int    `json:"work_count"`
	Works       []OpenLibraryBook `json:"works"`
}