package request

import (
	"time"
)

//MentoringRequestsResults represents the complete response when requesting mentoring requests
// at the Exercism API
type MentoringRequestsResults struct {
	MentoringRequests []MentoringRequest `json:"results"`
	Meta              Meta               `json:"meta"`
}

//Meta represents the meta data at the Exercism API
type Meta struct {
	CurrentPage int `json:"current_page"`
	TotalCount  int `json:"total_count"`

	TotalPages    int `json:"total_pages"`
	UnscopedTotal int `json:"unscoped_total"`
}

//MentoringRequest represents one mentoring request at the Exercism API
type MentoringRequest struct {
	UUID                   string      `json:"uuid"`
	TrackTitle             string      `json:"track_title"`
	ExerciseIconURL        string      `json:"exercise_icon_url"`
	ExerciseTitle          string      `json:"exercise_title"`
	StudentHandle          string      `json:"student_handle"`
	StudentAvatarURL       string      `json:"student_avatar_url"`
	UpdatedAt              time.Time   `json:"updated_at"`
	HaveMentoredPreviously bool        `json:"have_mentored_previously"`
	IsFavorited            bool        `json:"is_favorited"`
	Status                 interface{} `json:"status"`
	TooltipURL             string      `json:"tooltip_url"`
	URL                    string      `json:"url"`
}
