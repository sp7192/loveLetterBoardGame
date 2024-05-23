package card

type ICardEffect interface {
	Play(ownerId, targetId uint, data string) error
	Discard() error
}
