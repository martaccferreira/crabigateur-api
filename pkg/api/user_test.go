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

func (m *MockUserRepository) GetReviews(userId string, numReviews int, sort []api.SortOrder) ([]api.Card, error) {
	args := m.Called(userId, numReviews, sort)
	return args.Get(0).([]api.Card), args.Error(1)
}

func (m *MockUserRepository) InsertReviews(userId string, reviews []api.Review) ([]api.ReviewResult, error) {
	args := m.Called(userId, reviews)
	return args.Get(0).([]api.ReviewResult), args.Error(1)
}

func (m *MockUserRepository) UpdateReviews(userId string, reviews []api.Review) ([]api.ReviewResult, error) {
	args := m.Called(userId, reviews)
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
		expectErr  bool
	}{
		{
			name:       "Success",
			userId:     "123",
			numLessons: 5,
			mockResult: []api.Card{{CardId: 1, Word: "word1"}},
			mockError:  nil,
			expected:   []api.Card{{CardId: 1, Word: "word1"}},
			expectErr:  false,
		},
		{
			name:       "Repository Error",
			userId:     "123",
			numLessons: 5,
			mockResult: nil,
			mockError:  errors.New("repository error"),
			expected:   nil,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
            service := api.NewUserService(mockRepo)

			mockRepo.On("GetLessons", tt.userId, tt.numLessons).Return(tt.mockResult, tt.mockError)

			result, err := service.LessonCards(tt.userId, tt.numLessons)
			log.Println(err)

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

func TestUserService_ReviewCards(t *testing.T) {
	tests := []struct {
		name       string
		userId     string
		numReviews int
		sort       []api.SortOrder
		mockResult []api.Card
		mockError  error
		expected   []api.Card
		expectErr  bool
	}{
		{
			name:       "Success",
			userId:     "123",
			numReviews: 10,
			sort:       []api.SortOrder{api.DateAsc},
			mockResult: []api.Card{{CardId: 1, Word: "word1"}},
			mockError:  nil,
			expected:   []api.Card{{CardId: 1, Word: "word1"}},
			expectErr:  false,
		},
		{
			name:       "Repository Error",
			userId:     "123",
			numReviews: 10,
			sort:       []api.SortOrder{api.DateAsc},
			mockResult: nil,
			mockError:  errors.New("repository error"),
			expected:   nil,
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
            service := api.NewUserService(mockRepo)

			mockRepo.On("GetReviews", tt.userId, tt.numReviews, tt.sort).Return(tt.mockResult, tt.mockError)

			result, err := service.ReviewCards(tt.userId, tt.numReviews, tt.sort)

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
		reviews   []api.Review
		mockResult []api.ReviewResult
		mockError  error
		expected  []api.ReviewResult
		expectErr bool
	}{
		{
			name:      "Success",
			userId:    "123",
			reviews:   []api.Review{{CardId: 1, Success: new(bool)}},
			mockResult: []api.ReviewResult{{CardId: 1, Success: true}},
			mockError:  nil,
			expected:  []api.ReviewResult{{CardId: 1, Success: true}},
			expectErr: false,
		},
		{
			name:      "Repository Error",
			userId:    "123",
			reviews:   []api.Review{{CardId: 1, Success: new(bool)}},
			mockResult: nil,
			mockError:  errors.New("repository error"),
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := api.NewUserService(mockRepo)

			mockRepo.On("InsertReviews", tt.userId, tt.reviews).Return(tt.mockResult, tt.mockError)

			result, err := service.AddReviews(tt.userId, tt.reviews)

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
		reviews   []api.Review
		mockResult []api.ReviewResult
		mockError  error
		expected  []api.ReviewResult
		expectErr bool
	}{
		{
			name:      "Success",
			userId:    "123",
			reviews:   []api.Review{{CardId: 1, Success: new(bool)}},
			mockResult: []api.ReviewResult{{CardId: 1, Success: true}},
			mockError:  nil,
			expected:  []api.ReviewResult{{CardId: 1, Success: true}},
			expectErr: false,
		},
		{
			name:      "Repository Error",
			userId:    "123",
			reviews:   []api.Review{{CardId: 1, Success: new(bool)}},
			mockResult: nil,
			mockError:  errors.New("repository error"),
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			service := api.NewUserService(mockRepo)

			mockRepo.On("UpdateReviews", tt.userId, tt.reviews).Return(tt.mockResult, tt.mockError)

			result, err := service.UpdateReviews(tt.userId, tt.reviews)

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
