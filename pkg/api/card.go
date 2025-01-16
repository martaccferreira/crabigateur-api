package api

type CardService interface {
	GetCardById(id int) (Card, error)
}

type CardRepository interface {
	GetCard(id int) (Card, error)
}

type cardService struct {
	storage CardRepository
}

func NewCardService(cardRepo CardRepository) CardService {
	return &cardService{
		storage: cardRepo,
	}
}

func (u* cardService) GetCardById(id int) (Card, error) {
	card, err := u.storage.GetCard(id)
	if err != nil {
		return Card{}, err
	}
	return card, nil
}