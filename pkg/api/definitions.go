package api

import "time"

type Card struct {
	CardId int `json:"card_id"  binding:"numeric"`
	Level int `json:"level" binding:"required,numeric"`
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

type FormFlexion string

const (
	MascSing FormFlexion = "m.s."
	MascPlur FormFlexion = "m.p."
	FemSing FormFlexion = "f.s."
	FemPlur FormFlexion = "f.p."
)

type FormNumber string

const (
	Sing FormNumber = "singular"
	Plur FormNumber = "plural"
)

type FormGender string

const (
	Masc FormNumber = "m"
	Fem FormNumber = "f"
)

type WordTypes string

const (
	Regular WordTypes = "regular"
	Irregular WordTypes = "irregular"
	Verb WordTypes = "verb"
)

type UserPath struct {
	UserId string `uri:"user_id" binding:"required,numeric"`
}

type CardPath struct {
	CardId int `uri:"card_id" binding:"required,numeric"`
}

type QueryParams struct {
	FirstReview bool `form:"first_review" binding:"omitempty"`
	Sort []SortOrder `form:"sort" binding:"omitempty,sortable"`
	NumCards int	`form:"num_cards" binding:"omitempty,numeric,gte=0"`
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

type QuizList struct {
	CardIds []int `json:"card_ids" binding:"required"`
}

type QuizSummary struct {
	StageId int             `json:"stage_id"`
	Cards   []ReviewResult `json:"cards"`
}

