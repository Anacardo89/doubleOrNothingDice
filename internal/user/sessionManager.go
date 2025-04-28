package user

import (
	"context"
	"errors"
	"sync"

	"github.com/Anacardo89/doubleOrNothingDice/internal/db"
	"github.com/Anacardo89/doubleOrNothingDice/internal/game"
)

type Session struct {
	ClientID string
	Balance  int
	Game     *game.Game
}

type SessionManager struct {
	sessions map[string]*Session
	DB       *db.Manager
	mu       sync.RWMutex
}

func NewSessionManager(db *db.Manager) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		DB:       db,
	}
}

func (sm *SessionManager) Create(clientID string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	var user *db.User
	ctx := context.Background()
	user, err := sm.DB.GetUserByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	session := &Session{
		ClientID: clientID,
		Balance:  user.Balance,
	}
	sm.sessions[clientID] = session
	return session, nil
}

func (sm *SessionManager) Get(clientID string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.sessions[clientID]
	return s, ok
}

func (sm *SessionManager) Delete(clientID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, clientID)
}

func (sm *SessionManager) StartGame(clientID string, initialBet int) (*game.Game, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, ok := sm.sessions[clientID]
	if !ok {
		return nil, errors.New("session not found")
	}
	if session.Balance < initialBet {
		return nil, errors.New("insufficient balance")
	}
	session.Balance -= initialBet
	ctx := context.Background()
	err := sm.DB.UpdateUserBalance(ctx, session.ClientID, session.Balance)
	if err != nil {
		session.Balance += initialBet
		return nil, err
	}
	session.Game = game.NewGame(clientID, initialBet)
	dbgame := &db.Game{
		UserID:     session.Game.ClientID,
		InitialBet: session.Game.InitialBet,
		FinalBet:   session.Game.CurrentBet,
		TotalPlays: len(session.Game.Plays),
	}
	err = sm.DB.CreateGame(ctx, dbgame)
	if err != nil {
		return nil, err
	}
	session.Game.ID = dbgame.ID
	return session.Game, nil
}

func (sm *SessionManager) PlayRound(clientID string, betType string) (*game.Play, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session, ok := sm.sessions[clientID]
	if !ok || session.Game == nil {
		return nil, errors.New("no active game")
	}
	playResult, err := session.Game.Play(betType)
	if err != nil {
		return nil, err
	}
	if playResult.Outcome == game.OutcomeLose {
		ctx := context.Background()
		dbgame := &db.Game{
			ID:         session.Game.ID,
			FinalBet:   session.Game.CurrentBet,
			TotalPlays: len(session.Game.Plays),
		}
		if err := sm.DB.UpdateGame(ctx, dbgame); err != nil {
			return nil, err
		}
		for i, play := range session.Game.Plays {
			dbplay := &db.Play{
				GameID:     session.Game.ID,
				PlayNumber: i + 1,
				BetAmount:  play.BetAmount,
				PlayChoice: play.PlayChoice,
				DiceResult: play.DiceResult,
				Outcome:    play.Outcome,
			}
			if err := sm.DB.CreatePlay(ctx, dbplay); err != nil {
				return nil, err
			}
		}
	}
	return playResult, nil
}

func (sm *SessionManager) EndGame(clientID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session, ok := sm.sessions[clientID]
	if !ok || session.Game == nil {
		return errors.New("no active game to end")
	}
	session.Game.EndGame()
	session.Balance += session.Game.CurrentBet
	ctx := context.Background()
	dbgame := &db.Game{
		ID:         session.Game.ID,
		FinalBet:   session.Game.CurrentBet,
		TotalPlays: len(session.Game.Plays),
	}
	if err := sm.DB.UpdateGame(ctx, dbgame); err != nil {
		return err
	}
	for i, play := range session.Game.Plays {
		dbplay := &db.Play{
			GameID:     session.Game.ID,
			PlayNumber: i + 1,
			BetAmount:  play.BetAmount,
			PlayChoice: play.PlayChoice,
			DiceResult: play.DiceResult,
			Outcome:    play.Outcome,
		}
		if err := sm.DB.CreatePlay(ctx, dbplay); err != nil {
			return err
		}
	}
	if err := sm.DB.UpdateUserBalance(ctx, session.ClientID, session.Balance); err != nil {
		return err
	}
	session.Game = nil
	return nil
}
