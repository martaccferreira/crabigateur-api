package api_test

import (
	"crabigateur-api/pkg/api"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestSortOrder_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		sort     api.SortOrder
		expected bool
	}{
		{name: "Valid DateAsc", sort: api.DateAsc, expected: true},
		{name: "Valid DateDesc", sort: api.DateDesc, expected: true},
		{name: "Valid LevelAsc", sort: api.LevelAsc, expected: true},
		{name: "Valid LevelDesc", sort: api.LevelDesc, expected: true},
		{name: "Invalid SortOrder", sort: "invalid_sort", expected: false},
		{name: "Empty", sort: "", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.sort.IsValid()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestValidSortOrders(t *testing.T) {
	validate := validator.New()
	_ = validate.RegisterValidation("sortable", api.ValidSortOrders)

	tests := []struct {
		name       string
		orders     []api.SortOrder
		expected   bool
	}{
		{name: "Valid single sort", orders: []api.SortOrder{api.DateAsc}, expected: true},
		{name: "Valid multiple sorts", orders: []api.SortOrder{api.DateAsc, api.LevelDesc}, expected: true},
		{name: "Empty sort orders", orders: []api.SortOrder{}, expected: true},
		{name: "Invalid sort order", orders: []api.SortOrder{"invalid_sort"}, expected: false},
		{name: "Conflict between DateAsc and DateDesc", orders: []api.SortOrder{api.DateAsc, api.DateDesc}, expected: false},
		{name: "Conflict between LevelAsc and LevelDesc", orders: []api.SortOrder{api.LevelAsc, api.LevelDesc}, expected: false},
		{name: "Valid with multiple non-conflicting orders", orders: []api.SortOrder{api.DateAsc, api.LevelAsc}, expected: true},
		{name: "Conflicting orders in mixed list", orders: []api.SortOrder{api.DateAsc, api.LevelDesc, api.DateDesc}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.orders, "sortable")
			if tt.expected {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
