package repository_test

import (
	"crabigateur-api/pkg/api"
	"crabigateur-api/pkg/repository"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetLessons(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := repository.NewStorage(db)

	tests := []struct {
		name       string
		userId     string
		numLessons int
		mockSetup  func()
		expected   []api.Card
		expectErr  bool
	}{
		{
			name:       "Success",
			userId:     "123",
			numLessons: 2,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"card_id", "word", "translation", "word_type", "level",
					"tense", "forms", "irregular", "gender", "number", "form",
				}).
					AddRow(1, "chat", []byte(`["cat"]`), "regular", 1, nil, nil, nil, nil, nil, nil).
					AddRow(2, "chien", []byte(`["dog"]`), "regular", 1, nil, nil, nil, nil, nil, nil)
				mock.ExpectQuery(`WITH PendingLessonCardIds AS .*`).
					WithArgs("123").
					WillReturnRows(rows)
			},
			expected: []api.Card{
				{CardId: 1, Word: "chat", Translation: []string{"cat"}, WordType: "regular", Level: 1},
				{CardId: 2, Word: "chien", Translation: []string{"dog"}, WordType: "regular", Level: 1},
			},
			expectErr: false,
		},
		{
			name:       "Empty Result Set",
			userId:     "123",
			numLessons: 2,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"card_id", "word", "translation", "word_type", "level",
					"tense", "forms", "irregular", "gender", "number", "form",
				})
				mock.ExpectQuery(`WITH PendingLessonCardIds AS .*`).
					WithArgs("123").
					WillReturnRows(rows)
			},
			expected:  []api.Card{},
			expectErr: false,
		},
		{
			name:       "SQL Error",
			userId:     "123",
			numLessons: 2,
			mockSetup: func() {
				mock.ExpectQuery(`WITH PendingLessonCardIds AS .*`).
					WithArgs("123").
					WillReturnError(fmt.Errorf("query error"))
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name:       "Invalid Data in Columns",
			userId:     "123",
			numLessons: 2,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"card_id", "word", "translation", "word_type", "level",
					"tense", "forms", "irregular", "gender", "number", "form",
				}).
					AddRow("invalid_id", "chat", []byte(`["cat"]`), "regular", 1, nil, nil, nil, nil, nil, nil)
				mock.ExpectQuery(`WITH PendingLessonCardIds AS .*`).
					WithArgs("123").
					WillReturnRows(rows)
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			cards, err := storage.GetLessons(tt.userId, tt.numLessons)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, cards)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, cards)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetReviews(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := repository.NewStorage(db)

	tests := []struct {
		name       string
		userId     string
		numReviews int
		sort       []api.SortOrder
		mockSetup  func()
		expected   []api.Card
		expectErr  bool
	}{
		{
			name:       "Success",
			userId:     "123",
			numReviews: 2,
			sort:       []api.SortOrder{api.DateAsc, api.LevelDesc},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"card_id", "word", "translation", "word_type", "level",
					"tense", "forms", "irregular", "gender", "number", "form",
				}).
					AddRow(1, "chat", []byte(`["cat"]`), "regular", 1, nil, nil, nil, nil, nil, nil).
					AddRow(2, "chien", []byte(`["dog"]`), "regular", 1, nil, nil, nil, nil, nil, nil)
				mock.ExpectQuery(`WITH PendingReviews AS .*`).
					WithArgs("123").
					WillReturnRows(rows)
			},
			expected: []api.Card{
				{CardId: 1, Word: "chat", Translation: []string{"cat"}, WordType: "regular", Level: 1},
				{CardId: 2, Word: "chien", Translation: []string{"dog"}, WordType: "regular", Level: 1},
			},
			expectErr: false,
		},
		{
			name:       "Empty Result Set",
			userId:     "123",
			numReviews: 2,
			sort:       []api.SortOrder{api.DateAsc},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"card_id", "word", "translation", "word_type", "level",
					"tense", "forms", "irregular", "gender", "number", "form",
				})
				mock.ExpectQuery(`WITH PendingReviews AS .*`).
					WithArgs("123").
					WillReturnRows(rows)
			},
			expected:  []api.Card{},
			expectErr: false,
		},
		{
			name:       "SQL Error",
			userId:     "123",
			numReviews: 2,
			sort:       []api.SortOrder{api.DateAsc},
			mockSetup: func() {
				mock.ExpectQuery(`WITH PendingReviews AS .*`).
					WithArgs("123").
					WillReturnError(fmt.Errorf("query error"))
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name:       "Invalid Data in Columns",
			userId:     "123",
			numReviews: 2,
			sort:       []api.SortOrder{api.DateAsc},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"card_id", "word", "translation", "word_type", "level",
					"tense", "forms", "irregular", "gender", "number", "form",
				}).
					AddRow("invalid_id", "chat", []byte(`["cat"]`), "regular", 1, nil, nil, nil, nil, nil, nil)
				mock.ExpectQuery(`WITH PendingReviews AS .*`).
					WithArgs("123").
					WillReturnRows(rows)
			},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			cards, err := storage.GetReviews(tt.userId, tt.numReviews, tt.sort)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, cards)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, cards)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

