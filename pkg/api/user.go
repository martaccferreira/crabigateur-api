package api

type UserService interface {
	LessonCards(userId string, numLessons int) ([]Card, []int, error)
	ReviewCards(userId string, numReviews int, sort []SortOrder) ([]Card, []int, error)
	AddReviews(userId string, reviews []Review) ([]ReviewResult, error)
	UpdateReviews(userId string, reviews []Review) ([]ReviewResult, error)
	GetQuizSummary(userId string, cardIds []int) ([]ReviewResult, error)
}

type UserRepository interface {
	GetLessons(userId string, numLessons int) ([]Card, error)
	GetReviews(userId string, numReviews int, sort []SortOrder) ([]Card, error)
	InsertReview(userId string, review Review) (ReviewResult, error)
	UpdateReview(userId string, reviews Review) (ReviewResult, error)
	GetMostRecentReview(userId string, cardId int) (ReviewResult, error)
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

func (u* userService) ReviewCards(userId string, numReviews int, sort []SortOrder) ([]Card, []int, error) {
	reviews, err := u.storage.GetReviews(userId, numReviews, sort)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]int, len(reviews))
	for i, card := range reviews {
		ids[i] = card.CardId
	}
	return reviews, ids, nil
}

func (u* userService) AddReviews(userId string, reviews []Review) ([]ReviewResult, error) {
	results := []ReviewResult{}
	for _, review := range reviews {
		result, err := u.storage.InsertReview(userId, review)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (u* userService) UpdateReviews(userId string, reviews []Review) ([]ReviewResult, error) {
	results := []ReviewResult{}
	for _, review := range reviews {
		result, err := u.storage.UpdateReview(userId, review)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (u* userService) GetQuizSummary(userId string, cardIds []int) ([]ReviewResult, error) {
	results := []ReviewResult{}
	for _, cardId := range cardIds {
		review, err := u.storage.GetMostRecentReview(userId, cardId)
		if err != nil {
			return nil, err
		}
		results = append(results, review)
	}
	return results, nil
}
