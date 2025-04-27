package game

import (
	"errors"
	"math/rand"
	"time"
)

const (
	BetEven = "even"
	BetOdd  = "odd"

	OutcomeWin  = "win"
	OutcomeLose = "lose"
)

type Game struct {
	ID         string
	ClientID   string
	InitialBet int
	CurrentBet int
	Plays      []Play
	IsActive   bool
	rng        *rand.Rand
}

type Play struct {
	BetAmount  int
	PlayChoice string
	Outcome    string
	DiceResult int
}

// NewGame initializes a new game for a player.
func NewGame(clientID string, initialBet int) *Game {
	return &Game{
		ClientID:   clientID,
		InitialBet: initialBet,
		CurrentBet: initialBet,
		IsActive:   true,
		Plays:      []Play{},
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *Game) RollDice() int {
	return g.rng.Intn(6) + 1
}

func (g *Game) Play(playChoice string) (*Play, error) {
	if !g.IsActive {
		return nil, errors.New("game is not active")
	}
	diceResult := g.RollDice()
	outcome := OutcomeLose
	if (diceResult%2 == 0 && playChoice == BetEven) || (diceResult%2 != 0 && playChoice == BetOdd) {
		outcome = OutcomeWin
	}
	play := Play{
		BetAmount:  g.CurrentBet,
		PlayChoice: playChoice,
		Outcome:    outcome,
		DiceResult: diceResult,
	}
	g.Plays = append(g.Plays, play)
	if outcome == OutcomeLose {
		g.IsActive = false
		g.CurrentBet = 0
	}
	if g.IsActive {
		g.CurrentBet *= 2
	}
	return &play, nil
}

func (g *Game) EndGame() {
	g.IsActive = false
}
