package api

type CardService interface {
	GetCardById(id int) (Card, error)
	CreateCard(card Card) (Card, error)
	UpdateCard(cardId int, card Card) (Card, error)
	DeleteCard(cardId int) error
	SearchCards(query CardQueryParams) ([]Card, error)
}

type CardRepository interface {
	GetCard(id int) (Card, error)
	InsertCard(word string, translation []string, wordType string, gender string, level int) (int, error)
	UpdateCard(cardId int, word string, translation []string, wordType string, gender string, level int) error
	InsertOrUpdateConjugation(isUpdate bool, cardId int, tense string, forms []string, isIrregular bool) error
	InsertOrUpdateForm(isUpdate bool, cardId int, gender string, number string, form string) error
	DeleteCard(cardId int) error
	SearchCards(query CardQueryParams) ([]Card, error)
}

type cardService struct {
	storage CardRepository
}

func NewCardService(cardRepo CardRepository) CardService {
	return &cardService{
		storage: cardRepo,
	}
}

func (u *cardService) GetCardById(id int) (Card, error) {
	card, err := u.storage.GetCard(id)
	if err != nil {
		return Card{}, err
	}
	return card, nil
}

func (u *cardService) CreateCard(card Card) (Card, error) {
	cardId, err := u.storage.InsertCard(card.Word, card.Translation, card.WordType, card.Gender, card.Level)
	if err != nil {
		return Card{}, err
	}
	card.CardId = cardId

	err = u.processCardForms(false, card)
	if err != nil {
		return Card{}, err
	}

	return card, nil
}

func (u *cardService) UpdateCard(cardId int, card Card) (Card, error) {
	err := u.storage.UpdateCard(cardId, card.Word, card.Translation, card.WordType, card.Gender, card.Level)
	if err != nil {
		return Card{}, err
	}
	card.CardId = cardId // certifies cardId is set in case card payload doesn't include it

	err = u.processCardForms(true, card)
	if err != nil {
		return Card{}, err
	}

	return card, nil
}

func (u *cardService) processCardForms(isUpdate bool, card Card) error {
	for key, value := range card.Forms {
		var err error
		switch {
		case card.WordType == "verb":
			err = u.storage.InsertOrUpdateConjugation(isUpdate, card.CardId, key, value, card.IrregularVerb)
		case card.WordType == "irregular":
			err = u.insertIrregularForms(isUpdate, card.CardId, key, value)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *cardService) insertIrregularForms(isUpdate bool, cardId int, key string, value []string) error {
	var checkedValue string
	if len(value) == 0 {
		checkedValue = ""
	} else {
		checkedValue = value[0]
	}

	switch FormFlexion(key) {
	case MascSing:
		return u.storage.InsertOrUpdateForm(isUpdate, cardId, string(Masc), string(Sing), checkedValue)
	case MascPlur:
		return u.storage.InsertOrUpdateForm(isUpdate, cardId, string(Masc), string(Plur), checkedValue)
	case FemSing:
		return u.storage.InsertOrUpdateForm(isUpdate, cardId, string(Fem), string(Sing), checkedValue)
	case FemPlur:
		return u.storage.InsertOrUpdateForm(isUpdate, cardId, string(Fem), string(Plur), checkedValue)
	}
	return nil
}

func (u *cardService) DeleteCard(cardId int) error {
	err := u.storage.DeleteCard(cardId)
	if err != nil {
		return err
	}

	return nil
}

func (u *cardService) SearchCards(query CardQueryParams) ([]Card, error) {
	return u.storage.SearchCards(query)
}
