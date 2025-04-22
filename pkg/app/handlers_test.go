package app_test

import (
	"bytes"
	"crabigateur-api/pkg/api"
	"crabigateur-api/pkg/app"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

type fields struct {
	userService *MockService
	cardService *MockService
}

func (m *MockService) LessonCards(userId string, numLessons int) ([]api.Card, []int, error) {
	args := m.Called(userId, numLessons)
	return args.Get(0).([]api.Card), args.Get(1).([]int), args.Error(2)
}

func (m *MockService) ReviewCard(userId string, firstReview bool, sort []api.SortOrder) (api.Card, error) {
	args := m.Called(userId, firstReview, sort)
	return args.Get(0).(api.Card), args.Error(1)
}

func (m *MockService) AddReviews(userId string, cardIds []int) ([]api.ReviewResult, error) {
	args := m.Called(userId, cardIds)
	return args.Get(0).([]api.ReviewResult), args.Error(1)
}

func (m *MockService) UpdateReview(userId string, review api.Review) (api.ReviewResult, error) {
	args := m.Called(userId, review)
	return args.Get(0).(api.ReviewResult), args.Error(1)
}

func (m *MockService) GetQuizSummary(userId string, numCards int) ([]api.QuizSummary, error) {
	args := m.Called(userId, numCards)
	return args.Get(0).([]api.QuizSummary), args.Error(1)
}

func (m *MockService) GetCardById(id int) (api.Card, error) {
	args := m.Called(id)
	return args.Get(0).(api.Card), args.Error(1)
}

func (m *MockService) CreateCard(card api.Card) (api.Card, error) {
	args := m.Called(card)
	return args.Get(0).(api.Card), args.Error(1)
}

func (m *MockService) UpdateCard(cardId int, card api.Card) (api.Card, error) {
	args := m.Called(card)
	return args.Get(0).(api.Card), args.Error(1)
}

func (m *MockService) DeleteCard(cardId int) error {
	args := m.Called(cardId)
	return args.Error(0)
}

func (m *MockService) SearchCards(query api.CardQueryParams) ([]api.Card, error) {
	args := m.Called(query)
	return args.Get(0).([]api.Card), args.Error(1)
}

func (m *MockService) GetStats(userId string) (map[string]interface{}, error) {
	args := m.Called(userId)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

type args struct {
	request func() *http.Request
}

func TestHandlers(t *testing.T) {
	tests := []struct {
		name               string
		fields             fields
		args               args
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:   "ApiStatus - Success",
			fields: fields{},
			args: args{
				request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/v1/api/status", nil)
					return req
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"data":"crabigateur API running smoothly"}`,
		},
		{
			name: "GetUserLessons - Success",
			fields: fields{
				userService: func() *MockService {
					mockService := new(MockService)
					mockService.On("LessonCards", "123", 10).Return([]api.Card{
						{CardId: 1, Word: "example"},
					}, []int{1}, nil)
					return mockService
				}(),
				cardService: nil,
			},
			args: args{
				request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/v1/api/lessons/123?num_cards=10", nil)
					return req
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"data":{"card_ids":[1],"cards":[{"card_id":1,"level":0,"word_type":"","translations":null,"word":"example","gender":"","forms":null,"is_irregular_verb":false}],"total":1}}`,
		},
		{
			name: "GetUserLessons - Invalid user_id",
			fields: fields{
				userService: new(MockService),
				cardService: nil,
			},
			args: args{
				request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/v1/api/lessons/a?num_cards=10", nil)
					return req
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid user_id"}`,
		},
		{
			name: "GetUserReviews - Success",
			fields: fields{
				userService: func() *MockService {
					mockService := new(MockService)
					mockService.On("ReviewCard", "123", false, []api.SortOrder{api.DateAsc}).Return(api.Card{
						CardId: 1, Word: "reviewed_word",
					}, nil)
					return mockService
				}(),
				cardService: nil,
			},
			args: args{
				request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/v1/api/reviews/123?sort=date_asc", nil)
					return req
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"data":{"card_id":1,"level":0,"word_type":"","translations":null,"word":"reviewed_word","gender":"","forms":null,"is_irregular_verb":false}}`,
		},
		{
			name: "GetUserReviews - Invalid query parameters",
			fields: fields{
				userService: new(MockService),
				cardService: nil,
			},
			args: args{
				request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/v1/api/reviews/123?num_reviews=5&sort=date_asc&sort=date_desc", nil)
					return req
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid query parameters"}`,
		},
		{
			name: "PostUserReviews - Success",
			fields: fields{
				userService: func() *MockService {
					mockService := new(MockService)
					mockService.On("AddReviews", "123", []int{1}).Return([]api.ReviewResult{
						{CardId: 1, Success: true, StageId: "2"},
					}, nil)
					return mockService
				}(),
				cardService: nil,
			},
			args: args{
				request: func() *http.Request {
					cardIds := map[string]interface{}{
						"card_ids": []int{1},
					}
					body, _ := json.Marshal(cardIds)
					req, _ := http.NewRequest(http.MethodPost, "/v1/api/reviews/123", bytes.NewReader(body))
					req.Header.Set("Content-Type", "application/json")
					return req
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"data":[{"card_id":1,"card_word":"","success":true,"stage_id":"2"}]}`,
		},
		{
			name: "PostUserReviews - Invalid review format (missing param)",
			fields: fields{
				userService: new(MockService),
				cardService: nil,
			},
			args: args{
				request: func() *http.Request {
					cardIds := map[string]interface{}{}
					body, _ := json.Marshal(cardIds)
					req, _ := http.NewRequest(http.MethodPost, "/v1/api/reviews/123", bytes.NewReader(body))
					req.Header.Set("Content-Type", "application/json")
					return req
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid review format"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := tt.fields.userService
			mockCardService := tt.fields.cardService

			router := gin.Default()
			server := app.NewServer(router, mockUserService, mockCardService)

			router = server.Routes()
			server.RegisterValidators()

			// Perform the request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, tt.args.request())

			// Validate the response
			if w.Code != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", w.Code, tt.expectedStatusCode)
			}

			if tt.expectedBody != "" {
				if gotBody := w.Body.String(); gotBody != tt.expectedBody {
					t.Errorf("handler returned unexpected body: got %v want %v", gotBody, tt.expectedBody)
				}
			}
		})
	}
}
