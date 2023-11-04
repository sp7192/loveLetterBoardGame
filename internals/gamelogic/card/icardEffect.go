package card

type ICardEffect interface {
	Play(ownerId, targetId uint, data string)
	Discard()
}
