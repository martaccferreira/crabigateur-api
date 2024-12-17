package repository

import (
	"crabigateur-api/pkg/api"
	"database/sql"
	"fmt"
)

type Storage interface {
	GetLessons(userId string, numLessons int) ([]api.Card, error)
	GetReviews(userId string, numReviews int, sort []api.SortOrder) ([]api.Card, error)
	InsertReviews(userId string, reviews []api.Review) ([]api.ReviewResult, error)
	UpdateReviews(userId string, reviews []api.Review) ([]api.ReviewResult, error)
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

func (s *storage) InsertReviews(userId string, reviews []api.Review) ([]api.ReviewResult, error){
	results := []api.ReviewResult{}

	for _, review := range reviews {
		_, err := s.ReviewsInsert(userId, review)
		if err != nil {
			return nil, fmt.Errorf("storage - Insert into Reviews Query: %s", err)
		}

		row := s.UserCardStatusInsert(userId, review)
		
		var result api.ReviewResult
		err = row.Scan(&result.CardId, &result.CardWord, &result.Success, &result.StageId)
		if err != nil {
			return nil, fmt.Errorf("storage - InsertReviews: %s", err)
		}

		results = append(results, result)
	}
	return results, nil
}

func (s *storage) UpdateReviews(userId string, reviews []api.Review) ([]api.ReviewResult, error){
	results := []api.ReviewResult{}

	for _, review := range reviews {
		_, err := s.ReviewsInsert(userId, review)
		if err != nil {
			return nil, fmt.Errorf("storage - Insert into Reviews Query: %s", err)
		}

		row := s.UserCardStatusUpdate(userId, review)
		
		var result api.ReviewResult
		err = row.Scan(&result.CardId, &result.CardWord, &result.Success, &result.StageId)
		if err != nil {
			return nil, fmt.Errorf("storage - UpdateReviews: %s", err)
		}

		results = append(results, result)
	}
	return results, nil
}
