package repository

import (
	"crabigateur-api/pkg/api"
	"database/sql"
	"fmt"
)

type Storage interface {
	GetLessons(userId string, numLessons int) ([]api.Card, error)
	GetReviews(userId string, numReviews int, sort []api.SortOrder) ([]api.Card, error)
	InsertReview(userId string, review api.Review) (api.ReviewResult, error)
	UpdateReview(userId string, review api.Review) (api.ReviewResult, error)
	GetMostRecentReview(userId string, cardId int) (api.ReviewResult, error)
	GetCard(id int) (api.Card, error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) GetLessons(userId string, numLessons int) ([]api.Card, error) {
	rows, err := s.LessonsQuery(userId, numLessons) 
	if err != nil {
		return nil, fmt.Errorf("storage - Get Lesson Cards Query: %s", err)
	}
	defer rows.Close()

	cards, err := parseAllCardsFromQuery(rows)
	if err != nil {
		return nil, fmt.Errorf("storage - GetLessons: %s", err)
	}

	return cards, nil
}

func (s *storage) GetReviews(userId string, numReviews int, sort []api.SortOrder) ([]api.Card, error) {
	rows, err := s.ReviewsQuery(userId, numReviews, sort)
	if err != nil {
		return nil, fmt.Errorf("storage - Get Review Cards Query: %s", err)
	}
	defer rows.Close()

	cards, err := parseAllCardsFromQuery(rows)
	if err != nil {
		return nil, fmt.Errorf("storage - GetReviews: %s", err)
	}

	return cards, nil
}

func (s *storage) InsertReview(userId string, review api.Review) (api.ReviewResult, error){
	_, err := s.ReviewsInsert(userId, review)
	if err != nil {
		return api.ReviewResult{}, fmt.Errorf("storage - Insert into Reviews Query: %s", err)
	}

	row := s.UserCardStatusInsert(userId, review)
	
	var result api.ReviewResult
	err = row.Scan(&result.CardId, &result.CardWord, &result.Success, &result.StageId)
	if err != nil {
		return api.ReviewResult{}, fmt.Errorf("storage - InsertReviews: %s", err)
		}

	return result, nil
}

func (s *storage) UpdateReview(userId string, review api.Review) (api.ReviewResult, error){
	_, err := s.ReviewsInsert(userId, review)
	if err != nil {
		return api.ReviewResult{}, fmt.Errorf("storage - Insert into Reviews Query: %s", err)
	}

	row := s.UserCardStatusUpdate(userId, review)
	
	var result api.ReviewResult
	err = row.Scan(&result.CardId, &result.CardWord, &result.Success, &result.StageId)
	if err != nil {
		return api.ReviewResult{}, fmt.Errorf("storage - UpdateReview: %s", err)
	}
	return result, nil
}

func (s *storage) GetMostRecentReview(userId string, cardId int) (api.ReviewResult, error){
	row := s.MostRecentReviewQuery(userId, cardId)
	
	var result api.ReviewResult
	err := row.Scan(&result.CardId, &result.CardWord, &result.Success, &result.StageId)
	if err != nil {
		return api.ReviewResult{}, fmt.Errorf("storage - GetMostRecentReview: %s", err)
	}
	return result, nil
}

func (s *storage) GetCard(id int) (api.Card, error) {
	rows, err := s.CardQuery(id) 
	if err == sql.ErrNoRows {
		return api.Card{}, nil
	} else if err != nil {
		return api.Card{}, fmt.Errorf("storage - GetCard: %s", err)
	}
	defer rows.Close()

	result, err := parseAllCardsFromQuery(rows)
	if err != nil {
		return api.Card{}, err
	}

	return result[0], nil

}
