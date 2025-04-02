package api

import (
	"fmt"
	"sort"
	"strconv"
)

type UserService interface {
	LessonCards(userId string, numLessons int) ([]Card, []int, error)
	ReviewCard(userId string, firstReview bool, sort []SortOrder) (Card, error)
	AddReviews(userId string, cardId []int) ([]ReviewResult, error)
	UpdateReview(userId string, review Review) (ReviewResult, error)
	GetQuizSummary(userId string, numCards int) ([]QuizSummary, error)
}

type UserRepository interface {
	GetLessons(userId string, numLessons int) ([]Card, error)
	GetReview(userId string, firstReview bool, sort []SortOrder) ([]Card, error)
	InsertReview(userId string, cardId int) (ReviewResult, error)
	UpdateReview(userId string, reviews Review) (ReviewResult, error)
	GetMostRecentReviews(userId string, numCards int) ([]ReviewResult, error)
}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}

func (u* userService) LessonCards(userId string, numLessons int) ([]Card, []int, error) {
	cards, err := u.storage.GetLessons(userId, numLessons)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]int, len(cards))
	for i, card := range cards {
		ids[i] = card.CardId
	}
	return cards, ids, nil
}

func (u* userService) ReviewCard(userId string, firstReview bool, sort []SortOrder) (Card, error) {
	reviews, err := u.storage.GetReview(userId, firstReview, sort)
	if err != nil {
		return Card{}, err
	}
	if len(reviews) == 0 {
		return Card{}, nil
	} 
	return reviews[0], nil
}

func (u* userService) AddReviews(userId string, cardIds []int) ([]ReviewResult, error) {
	results := []ReviewResult{}
	for _, cardId := range cardIds {
		result, err := u.storage.InsertReview(userId, cardId)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (u* userService) UpdateReview(userId string, review Review) (ReviewResult, error) {
	result, err := u.storage.UpdateReview(userId, review)
	if err != nil {
		return ReviewResult{}, err
	}
		
	return result, nil
}

func (u* userService) GetQuizSummary(userId string, numCards int) ([]QuizSummary, error) {
	reviews, err := u.storage.GetMostRecentReviews(userId, numCards)
	if err != nil {
		return nil, err
	}
	
	// Map to group reviews by stage_id
	groupedReviews := make(map[int][]ReviewResult)

	for _, review := range reviews {
		stageId, err := strconv.Atoi(review.StageId) // Convert stage_id from string to int
		if err != nil {
			return nil, fmt.Errorf("invalid stage_id format: %v", err)
		}
		groupedReviews[stageId] = append(groupedReviews[stageId], review)
	}

	// Convert the map into a slice
	var summary []QuizSummary
	for stageId, cards := range groupedReviews {
		summary = append(summary, QuizSummary{
			StageId: stageId,
			Cards:   cards,
		})
	}

	sort.Slice(summary, func(i, j int) bool {
		return summary[i].StageId < summary[j].StageId
	})

	return summary, nil
}
