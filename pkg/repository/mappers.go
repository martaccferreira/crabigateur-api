package repository

import (
	"crabigateur-api/pkg/api"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type Card struct {
	CardId int
	Level int
	WordType string
	Translation json.RawMessage
	Word string
	Gender sql.NullString
}

type Verb struct {
	Tense sql.NullString
	Forms *json.RawMessage
	Irregular sql.NullBool
}

type Form struct {
	Gender sql.NullString
	Number sql.NullString
	Form sql.NullString
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return "" // Default value if NULL
}

func nullBoolToBool(ns sql.NullBool) bool {
	if ns.Valid {
		return ns.Bool
	}
	return false 
}

func getFormKey(f Form) string {
	number := nullStringToString(f.Number)
	gender := nullStringToString(f.Gender)

	if(number == "singular"){
		return fmt.Sprintf("%s.s.", gender)
	}
	return fmt.Sprintf("%s.p.", gender)
}

func getConjugationKey(v Verb) string {
	return nullStringToString(v.Tense)
}

func getConjugationForms(f *json.RawMessage) ([]string, error) {
	if f != nil {
		var unmarshalledForms []string
		err := json.Unmarshal(*f, &unmarshalledForms)
		if err != nil {
			return nil, err
		}
		return unmarshalledForms, nil
	} 
	return []string{}, nil
}

func addFormToCard(card api.Card, form Form) {
	card.Forms[getFormKey(form)] = []string{nullStringToString(form.Form)}
}

func addConjugationToCard(card api.Card, verb Verb) error {
	var err error
	card.Forms[getConjugationKey(verb)], err = getConjugationForms(verb.Forms)
	if err != nil {
		return err
	}
	return nil
}

func parseNewCardFromRow(card Card, verb Verb, form Form) (api.Card, error) {
	var unmarshalledTranslations []string
	err := json.Unmarshal(card.Translation, &unmarshalledTranslations)
	if err != nil {
		return api.Card{}, err
	}

	newCard := api.Card{
		CardId: card.CardId,
		Level: card.Level,
		WordType: card.WordType,
		Translation: unmarshalledTranslations,
	}

	switch card.WordType {
	case "regular":
			newCard.Word = card.Word
			newCard.Gender = nullStringToString(card.Gender)
	case "irregular":
			newCard.Forms = make(map[string][]string)
			newCard.Forms["m.s."] = []string{card.Word}
			addFormToCard(newCard, form)
	case "verb":
			newCard.Word = card.Word
			newCard.Forms = make(map[string][]string)
			err = addConjugationToCard(newCard, verb)
			if err != nil {
				return api.Card{}, err
			}
			newCard.IrregularVerb = nullBoolToBool(verb.Irregular)
		
	}
	return newCard, nil
}

func addMoreFormsToExistingCard(card api.Card, verb Verb, form Form) error {
	switch card.WordType {
	case "irregular":
		addFormToCard(card, form)
	case "verb":
		return addConjugationToCard(card, verb)
	}
	return nil
}

func parseAllCardsFromQuery(rows *sql.Rows) ([]api.Card, error) {
	var newCard api.Card
	cards := []api.Card{}
	lastCardId := 0 	// impossible card id to make sure uninitialized newCard isn't added to cards
	hasRows := false 	// flag that ensures cards are only appended if query has results
	
	for rows.Next() {
		hasRows = true
		var card Card
		var verb Verb
		var form Form

		err := rows.Scan(&card.CardId, &card.Word, &card.Translation, &card.WordType, &card.Level, &verb.Tense, &verb.Forms, &verb.Irregular, &form.Gender, &form.Number, &form.Form)
		if err != nil {
			return nil, err
		}

		if(card.CardId == lastCardId) {
			if err = addMoreFormsToExistingCard(newCard, verb, form); err != nil {
				log.Printf("storage - Error unmarshalling JSONB: %v", err)
				return nil, err
			}
			continue
		} 
		if lastCardId != 0 {
			cards = append(cards, newCard)
		}

		if newCard, err = parseNewCardFromRow(card, verb, form); err != nil {
			log.Printf("storage - Error unmarshalling JSONB: %v", err)
			return nil, err
		}

		lastCardId = card.CardId
	}

	if hasRows {
		cards = append(cards, newCard)
	}

	return cards, nil
}