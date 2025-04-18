package repository

import (
	"crabigateur-api/pkg/api"
	"database/sql"
	"fmt"
	"strings"
)

type Storage interface {
	GetLessons(userId string, numLessons int) ([]api.Card, error)
	GetReview(userId string, firstReview bool, sort []api.SortOrder) ([]api.Card, error)
	InsertReview(userId string, cardId int) (api.ReviewResult, error)
	UpdateReview(userId string, review api.Review) (api.ReviewResult, error)
	GetMostRecentReviews(userId string, numCards int) ([]api.ReviewResult, error)
	GetCard(id int) (api.Card, error)
	InsertCard(word string, translation []string, wordType string, gender string, level int) (int, error)
	UpdateCard(cardId int, word string, translation []string, wordType string, gender string, level int) error
	InsertOrUpdateConjugation(isUpdate bool, cardId int, tense string, forms []string, isIrregular bool) error
	InsertOrUpdateForm(isUpdate bool, cardId int, gender string, number string, form string) error
	DeleteCard(cardId int) error
	SearchCards(query api.CardQueryParams) ([]api.Card, error)
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

func (s *storage) InsertReview(userId string, cardId int) (api.ReviewResult, error) {
	row := s.UserCardStatusInsert(userId, cardId)

	var result api.ReviewResult
	err := row.Scan(&result.CardId, &result.CardWord, &result.Success, &result.StageId)
	if err != nil {
		return api.ReviewResult{}, fmt.Errorf("storage - InsertReview: %s", err)
	}

	return result, nil
}

func (s *storage) UpdateReview(userId string, review api.Review) (api.ReviewResult, error) {
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

func (s *storage) GetMostRecentReviews(userId string, numCards int) ([]api.ReviewResult, error) {
	rows, err := s.MostRecentReviewsQuery(userId, numCards)
	if err != nil {
		return nil, fmt.Errorf("storage - Get Most Recent Reviews Query: %s", err)
	}
	defer rows.Close()

	var result []api.ReviewResult
	for rows.Next() {
		var review api.ReviewResult
		err = rows.Scan(&review.CardId, &review.CardWord, &review.Success, &review.StageId)
		if err != nil {
			return nil, fmt.Errorf("storage - GetMostRecentReviews: %s", err)
		}
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

func (s *storage) InsertCard(word string, translation []string, wordType string, gender string, level int) (int, error) {
	row, err := s.CardsInsert(word, translation, wordType, gender, level)
	if err != nil {
		return 0, fmt.Errorf("storage - CardsInsert: %s", err)
	}

	var cardId int
	err = row.Scan(&cardId)
	if err != nil {
		if strings.Contains(err.Error(), "unique_word_gender") {
			return 0, fmt.Errorf("storage - InsertCard: duplicate card: %s", err)
		}
		return 0, fmt.Errorf("storage - InsertCard: %s", err)
	}
	return cardId, nil
}

func (s *storage) UpdateCard(cardId int, word string, translation []string, wordType string, gender string, level int) error {
	_, err := s.CardsUpdate(cardId, word, translation, wordType, gender, level)
	if err != nil {
		return fmt.Errorf("storage - InsertOrUpdateCard: %s", err)
	}
	return nil
}

func (s *storage) InsertOrUpdateConjugation(isUpdate bool, cardId int, tense string, forms []string, isIrregular bool) error {
	var err error

	if isUpdate {
		_, err = s.ConjugationsUpdate(cardId, tense, forms, isIrregular)
	} else {
		_, err = s.ConjugationsInsert(cardId, tense, forms, isIrregular)
	}
	if err != nil {
		return fmt.Errorf("storage - InsertOrUpdateConjugation: %s", err)
	}

	return nil
}

func (s *storage) InsertOrUpdateForm(isUpdate bool, cardId int, gender string, number string, form string) error {
	var err error

	if isUpdate {
		_, err = s.FormsUpdate(cardId, gender, number, form)
	} else {
		_, err = s.FormsInsert(cardId, gender, number, form)
	}
	if err != nil {
		return fmt.Errorf("storage - InsertOrUpdateForm: %s", err)
	}
	return nil
}

func (s *storage) DeleteCard(cardId int) error {
	_, err := s.CardsDelete(cardId)
	if err != nil {
		return fmt.Errorf("storage - DeleteCard: %s", err)
	}
	return nil
}

func (s *storage) SearchCards(query api.CardQueryParams) ([]api.Card, error) {
	rows, err := s.SearchCardsQuery(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []api.Card
	for rows.Next() {
		var card Card

		err := rows.Scan(&card.CardId, &card.Word, &card.Translation, &card.WordType, &card.Gender, &card.Level)
		if err != nil {
			return nil, err
		}
		newCard, err := parseCardNoForms(card)
		if err != nil {
			return nil, err
		}
		cards = append(cards, newCard)
	}

	return cards, nil
}
