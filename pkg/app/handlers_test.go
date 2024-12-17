package app_test

import (
	"bytes"
	"crabigateur-api/pkg/api"
	"crabigateur-api/pkg/app"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

type fields struct {
	userService *MockUserService
}

func (m *MockUserService) LessonCards(userId string, numLessons int) ([]api.Card, error) {
	args := m.Called(userId, numLessons)
	return args.Get(0).([]api.Card), args.Error(1)
}

func (m *MockUserService) ReviewCards(userId string, numReviews int, sort []api.SortOrder) ([]api.Card, error) {
	args := m.Called(userId, numReviews, sort)
	return args.Get(0).([]api.Card), args.Error(1)
}

func (m *MockUserService) AddReviews(userId string, reviews []api.Review) ([]api.ReviewResult, error) {
	args := m.Called(userId, reviews)
	return args.Get(0).([]api.ReviewResult), args.Error(1)
}

func (m *MockUserService) UpdateReviews(userId string, reviews []api.Review) ([]api.ReviewResult, error) {
	args := m.Called(userId, reviews)
	return args.Get(0).([]api.ReviewResult), args.Error(1)
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
			name: "ApiStatus - Success",
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
				userService: func() *MockUserService {
					mockService := new(MockUserService)
					mockService.On("LessonCards", "123", 10).Return([]api.Card{
						{CardId: 1, Word: "example"},
					}, nil)
					return mockService
				}(),
			},
			args: args{
				request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/v1/api/lessons/123?num_lessons=10", nil)
					return req
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"data":{"cards":[{"card_id":1,"level":0,"word_type":"","translations":null,"word":"example","gender":"","forms":null,"is_irregular_verb":false}],"total":1}}`,
		},
		{
			name: "GetUserLessons - Invalid user_id",
			fields: fields{
				userService: new(MockUserService),
			},
			args: args{
				request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/v1/api/lessons/a?num_lessons=10", nil)
					return req
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid user_id"}`,
		},
		{
			name: "GetUserReviews - Success",
			fields: fields{
				userService: func() *MockUserService {
					mockService := new(MockUserService)
					mockService.On("ReviewCards", "123", 5, []api.SortOrder{api.DateAsc}).Return([]api.Card{
						{CardId: 1, Word: "reviewed_word"},
					}, nil)
					return mockService
				}(),
			},
			args: args{
				request: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/v1/api/reviews/123?num_reviews=5&sort=date_asc", nil)
					return req
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"data":{"cards":[{"card_id":1,"level":0,"word_type":"","translations":null,"word":"reviewed_word","gender":"","forms":null,"is_irregular_verb":false}],"total":1}}`,
		},
		{
			name: "GetUserReviews - Invalid query parameters",
			fields: fields{
				userService: new(MockUserService),
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
				userService: func() *MockUserService {
					mockService := new(MockUserService)
					mockService.On("AddReviews", "123", mock.Anything).Return([]api.ReviewResult{
						{CardId: 1, Success: true, StageId: "2"},
					}, nil)
					return mockService
				}(),
			},
			args: args{
				request: func() *http.Request {
					reviews := api.Reviews{
						Reviews: []api.Review{
							{
								CardId:         1,
								ReviewDate:     time.Now(),
								Success:        func(b bool) *bool { return &b }(true),
								IncorrectCount: func(i int) *int { return &i }(0),
							},
						},
					}
					body, _ := json.Marshal(reviews)
					req, _ := http.NewRequest(http.MethodPost, "/v1/api/lessons/123", bytes.NewReader(body))
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
				userService: new(MockUserService),
			},
			args: args{
				request: func() *http.Request {
					reviews := api.Reviews{
						Reviews: []api.Review{
							{
								CardId:         1,
								ReviewDate:     time.Now(),
								Success:        func(b bool) *bool { return &b }(true),
							},
						},
					}
					body, _ := json.Marshal(reviews)
					req, _ := http.NewRequest(http.MethodPost, "/v1/api/lessons/123", bytes.NewReader(body))
					req.Header.Set("Content-Type", "application/json")
					return req
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid review format"}`,
		},
		{
			name: "PostUserReviews - Invalid review format (wrong date format)",
			fields: fields{
				userService: new(MockUserService),
			},
			args: args{
				request: func() *http.Request {
					reviews := map[string]interface{}{
						"reviews": []map[string]interface{}{
							{
								"card_id":         1,
								"review_date":     "2024-12-11",
								"success":        func(b bool) *bool { return &b }(true),
								"incorrect_count": func(i int) *int { return &i }(0),
							},
						},
					}
					body, _ := json.Marshal(reviews)
					req, _ := http.NewRequest(http.MethodPost, "/v1/api/lessons/123", bytes.NewReader(body))
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
			mockService := tt.fields.userService
			
			router := gin.Default()
			server := app.NewServer(router, mockService)

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
