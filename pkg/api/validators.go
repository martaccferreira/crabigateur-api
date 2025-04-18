package api

import (
	"fmt"
	"reflect"

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

func (c Card) IsEmpty() bool {
 return reflect.ValueOf(c).IsZero()
}

func CardStructValidation(sl validator.StructLevel) {
	card := sl.Current().Interface().(Card)

	// Validate WordType
	validTypes := map[string]bool{
		string(Regular):   true,
		string(Irregular): true,
		string(Verb):      true,
	}
	if _, ok := validTypes[card.WordType]; !ok {
		sl.ReportError(card.WordType, "WordType", "word_type", "wordtypevalid", "")
	}

	// Validate Gender
	if card.Gender != "" && card.Gender != string(Masc) && card.Gender != string(Fem) {
		sl.ReportError(card.Gender, "Gender", "gender", "gendervalid", "")
	}

	// Validate Forms only if word_type is "irregular" and Forms is not empty
	if card.WordType == string(Irregular) && len(card.Forms) > 0 {
		validForms := map[string]bool{
			string(MascSing): true,
			string(MascPlur): true,
			string(FemSing):  true,
			string(FemPlur):  true,
		}
		for key := range card.Forms {
			if !validForms[key] {
				sl.ReportError(key, "Forms", "forms", "formflexionvalid", fmt.Sprintf("invalid key: %s", key))
			}
		}
	}
}
