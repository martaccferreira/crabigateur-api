package api_test

import (
	"crabigateur-api/pkg/api"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetLessons(userId string, numLessons int) ([]api.Card, error) {
	args := m.Called(userId, numLessons)
	return args.Get(0).([]api.Card), args.Error(1)
}

func (m *MockUserRepository) GetReview(userId string, firstReview bool, sort []api.SortOrder) ([]api.Card, error) {
	args := m.Called(userId, firstReview, sort)
	return args.Get(0).([]api.Card), args.Error(1)
}

func (m *MockUserRepository) InsertReview(userId string, review api.Review) (api.ReviewResult, error) {
	args := m.Called(userId, review)
	return args.Get(0).(api.ReviewResult), args.Error(1)
}

func (m *MockUserRepository) UpdateReview(userId string, review api.Review) (api.ReviewResult, error) {
	args := m.Called(userId, review)
	return args.Get(0).(api.ReviewResult), args.Error(1)
}

func (m *MockUserRepository) GetMostRecentReviews(userId string, cardId int) ([]api.ReviewResult, error) {
	args := m.Called(userId, cardId)
	return args.Get(0).([]api.ReviewResult), args.Error(1)
}

func TestUserService_LessonCards(t *testing.T) {
	tests := []struct {
		name       string
		userId     string
		numLessons int
		mockResult []api.Card
		mockError  error
		expected   []api.Card
		expectedId []int
		expectErr  bool
	}{
		{
			name:       "Success",
			userId:     "123",
			numLessons: 5,
			mockResult: []api.Card{{CardId: 1, Word: "word1"}},
			mockError:  nil,
			expected:   []api.Card{{CardId: 1, Word: "word1"}},
			expectedId:   []int{1},
			expectErr:  false,
		},
		{
			name:       "Repository Error",
			userId:     "123",
			numLessons: 5,
			mockResult: nil,
			mockError:  errors.New("repository error"),
			expected:   nil,
			expectedId: nil,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
            service := api.NewUserService(mockRepo)

			mockRepo.On("GetLessons", tt.userId, tt.numLessons).Return(tt.mockResult, tt.mockError)

			result, resultId, err := service.LessonCards(tt.userId, tt.numLessons)
			log.Println(err)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
				assert.Equal(t, tt.expectedId, resultId)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_ReviewCards(t *testing.T) {
	tests := []struct {
		name       string
		userId     string
		firstReview bool
		sort       []api.SortOrder
		mockResult []api.Card
		mockError  error
		expected   api.Card
		expectErr  bool
	}{
		{
			name:       "Success",
			userId:     "123",
			firstReview: true,
			sort:       []api.SortOrder{api.DateAsc},
			mockResult: []api.Card{{CardId: 1, Word: "word1"}},
			mockError:  nil,
			expected:   api.Card{CardId: 1, Word: "word1"},
			expectErr:  false,
		},
		{
			name:       "Repository Error",
			userId:     "123",
			firstReview: true,
			sort:       []api.SortOrder{api.DateAsc},
			mockResult: nil,
			mockError:  errors.New("repository error"),
			expected:   api.Card{},
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
            service := api.NewUserService(mockRepo)

			mockRepo.On("GetReview", tt.userId, tt.firstReview, tt.sort).Return(tt.mockResult, tt.mockError)

			result, err := service.ReviewCard(tt.userId, tt.firstReview, tt.sort)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_AddReviews(t *testing.T) {
	tests := []struct {
		name      string
		userId    string
		review   api.Review
		mockResult api.ReviewResult
		mockError  error
		expected  api.ReviewResult
		expectErr bool
	}{
		{
			name:      "Success",
			userId:    "123",
			review:   api.Review{CardId: 1, Success: new(bool)},
			mockResult: api.ReviewResult{CardId: 1, Success: true},
			mockError:  nil,
			expected:  api.ReviewResult{CardId: 1, Success: true},
			expectErr: false,
		},
		{
			name:      "Repository Error",
			userId:    "123",
			review:   api.Review{CardId: 1, Success: new(bool)},
			mockResult: api.ReviewResult{},
			mockError:  errors.New("repository error"),
			expected:  api.ReviewResult{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := api.NewUserService(mockRepo)

			mockRepo.On("InsertReview", tt.userId, tt.review).Return(tt.mockResult, tt.mockError)

			result, err := service.AddReview(tt.userId, tt.review)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateReviews(t *testing.T) {
	tests := []struct {
		name      string
		userId    string
		review   api.Review
		mockResult api.ReviewResult
		mockError  error
		expected  api.ReviewResult
		expectErr bool
	}{
		{
			name:      "Success",
			userId:    "123",
			review:   api.Review{CardId: 1, Success: new(bool)},
			mockResult: api.ReviewResult{CardId: 1, Success: true},
			mockError:  nil,
			expected:  api.ReviewResult{CardId: 1, Success: true},
			expectErr: false,
		},
		{
			name:      "Repository Error",
			userId:    "123",
			review:   api.Review{CardId: 1, Success: new(bool)},
			mockResult: api.ReviewResult{},
			mockError:  errors.New("repository error"),
			expected:  api.ReviewResult{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := api.NewUserService(mockRepo)

			mockRepo.On("UpdateReview", tt.userId, tt.review).Return(tt.mockResult, tt.mockError)

			result, err := service.UpdateReview(tt.userId, tt.review)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
