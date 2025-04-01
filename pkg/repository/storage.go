package repository

import (
	"crabigateur-api/pkg/api"
	"database/sql"
	"fmt"
)

type Storage interface {
	GetLessons(userId string, numLessons int) ([]api.Card, error)
	GetReview(userId string, firstReview bool, sort []api.SortOrder) ([]api.Card, error)
	InsertReview(userId string, review api.Review) (api.ReviewResult, error)
	UpdateReview(userId string, review api.Review) (api.ReviewResult, error)
	GetMostRecentReviews(userId string, numCards int) ([]api.ReviewResult, error)
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

func (s *storage) GetReview(userId string, firstReview bool, sort []api.SortOrder) ([]api.Card, error) {
	rows, err := s.ReviewQuery(userId, firstReview, sort)
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
		return api.ReviewResult{}, fmt.Errorf("storage - Insert into Review Query: %s", err)
	}

	row := s.UserCardStatusInsert(userId, review)
	
	var result api.ReviewResult
	err = row.Scan(&result.CardId, &result.CardWord, &result.Success, &result.StageId)
	if err != nil {
		return api.ReviewResult{}, fmt.Errorf("storage - InsertReview: %s", err)
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

func (s *storage) GetMostRecentReviews(userId string, numCards int) ([]api.ReviewResult, error){
	rows, err := s.MostRecentReviewsQuery(userId, numCards)
	if err != nil {
		return nil, fmt.Errorf("storage - Get Most Recent Reviews Query: %s", err)
	}
	defer rows.Close()
	fmt.Print(rows)
	
	var result []api.ReviewResult
	for rows.Next() {
		var review api.ReviewResult
		err = rows.Scan(&review.CardId, &review.CardWord, &review.Success, &review.StageId)
		if err != nil {
			return nil, fmt.Errorf("storage - GetMostRecentReviews: %s", err)
		}
		fmt.Print(review)
		result = append(result, review)
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
