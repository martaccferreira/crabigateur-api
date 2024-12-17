package api

import (
	"github.com/go-playground/validator/v10"
)

func (s SortOrder) IsValid() bool {
	switch SortOrder(s) { 
	case DateAsc, DateDesc, LevelAsc, LevelDesc:
		return true
	default:
		return false
	}
}

var ValidSortOrders validator.Func = func(fl validator.FieldLevel) bool {
	orders, ok := fl.Field().Interface().([]SortOrder)
	if !ok {
		return false
	}
	
	conflicts := map[SortOrder]SortOrder{
		DateAsc:  DateDesc,
		DateDesc: DateAsc,
		LevelAsc: LevelDesc,
		LevelDesc: LevelAsc,
	}

	seen := make(map[SortOrder]bool)

	for _, value := range orders {
		order := SortOrder(value)
		if !order.IsValid() {
			return false
		}
		if conflict, exists := conflicts[order]; exists {
			if seen[conflict] {
				return false
			}
		}
		seen[order] = true
	}
	return true
}