package logic

import "loveLetterClient/internals/models"

type GameLogicBuilder struct {
	actions []func(*GameLogic)
}

func (b *GameLogicBuilder) PlayerIds(ids []uint) *GameLogicBuilder {
	b.actions = append(b.actions, func(g *GameLogic) {
		g.playersIdInGame = ids
	})
	return b
}

func (b *GameLogicBuilder) OwnHand(hand models.Hand) *GameLogicBuilder {
	b.actions = append(b.actions, func(g *GameLogic) {
		g.OwnHand = hand
	})
	return b
}

func (b *GameLogicBuilder) Build() *GameLogic {
	g := GameLogic{}
	for _, action := range b.actions {
		action(&g)
	}
	return &g
}
