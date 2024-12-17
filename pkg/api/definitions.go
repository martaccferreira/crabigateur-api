package api

import "time"

type Card struct {
	CardId int `json:"card_id"`
	Level int `json:"level"`
	WordType string `json:"word_type"`
	Translation []string `json:"translations"`
	Word string `json:"word"`
	Gender string `json:"gender"`
	Forms map[string][]string `json:"forms"`
	IrregularVerb bool `json:"is_irregular_verb"`
}

type SortOrder string

const (
	DateAsc SortOrder = "date_asc"
	DateDesc SortOrder = "date_desc"
	LevelAsc SortOrder = "level_asc"
	LevelDesc SortOrder = "level_desc"
)

type PathParams struct {
	UserId string `uri:"user_id" binding:"required,numeric"`
}

type QueryParams struct {
	NumLessons int `form:"num_lessons" binding:"omitempty,numeric,gte=0"`
	NumReviews int `form:"num_reviews" binding:"omitempty,numeric,gte=0"`
	Sort []SortOrder `form:"sort" binding:"omitempty,sortable"`
}

type Review struct {
	CardId int `json:"card_id" binding:"required"`
	ReviewDate time.Time `json:"review_date" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	Success *bool `json:"success" binding:"required"`
	IncorrectCount *int `json:"incorrect_count" binding:"required,gte=0"`
}

type Reviews struct {
	Reviews []Review `json:"reviews" binding:"dive"`
}

type ReviewResult struct {
	CardId int `json:"card_id"`
	CardWord string `json:"card_word"`
	Success bool `json:"success"`
	StageId string `json:"stage_id"`
}
