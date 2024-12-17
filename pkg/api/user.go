package api

type UserService interface {
	LessonCards(userId string, numLessons int) ([]Card, error)
	ReviewCards(userId string, numReviews int, sort []SortOrder) ([]Card, error)
	AddReviews(userId string, reviews []Review) ([]ReviewResult, error)
	UpdateReviews(userId string, reviews []Review) ([]ReviewResult, error)
}

type UserRepository interface {
	GetLessons(userId string, numLessons int) ([]Card, error)
	GetReviews(userId string, numReviews int, sort []SortOrder) ([]Card, error)
	InsertReviews(userId string, reviews []Review) ([]ReviewResult, error)
	UpdateReviews(userId string, reviews []Review) ([]ReviewResult, error)
}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}

func (u* userService) LessonCards(userId string, numLessons int) ([]Card, error) {
	cards, err := u.storage.GetLessons(userId, numLessons)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (u* userService) ReviewCards(userId string, numReviews int, sort []SortOrder) ([]Card, error) {
	reviews, err := u.storage.GetReviews(userId, numReviews, sort)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (u* userService) AddReviews(userId string, reviews []Review) ([]ReviewResult, error) {
	lessonResults, err := u.storage.InsertReviews(userId, reviews)
	if err != nil {
		return nil, err
	}
	return lessonResults, nil
}

func (u* userService) UpdateReviews(userId string, reviews []Review) ([]ReviewResult, error) {
	reviewResults, err := u.storage.UpdateReviews(userId, reviews)
	if err != nil {
		return nil, err
	}
	return reviewResults, nil
}
