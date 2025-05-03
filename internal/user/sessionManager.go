package user

import (
	"context"
	"errors"
	"sync"

	"github.com/Anacardo89/doubleOrNothingDice/internal/db"
	"github.com/Anacardo89/doubleOrNothingDice/internal/game"
	"github.com/Anacardo89/doubleOrNothingDice/internal/redis"
)

type Session struct {
	ClientID string
	Balance  int
	Game     *game.Game
}

type SessionManager struct {
	sessions map[string]*Session
	DB       *db.Manager
	Redis    *redis.Manager
	mu       sync.RWMutex
}

func NewSessionManager(db *db.Manager, redisManager *redis.Manager) *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		DB:       db,
		Redis:    redisManager,
	}
}

func (sm *SessionManager) Create(clientID string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
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
		TotalPlays: 0,
	}
	err = sm.DB.CreateGame(ctx, dbgame)
	if err != nil {
		return nil, err
	}
	session.Game.ID = dbgame.ID
	err = sm.Redis.SaveGame(ctx, session.Game)
	if err != nil {
		return nil, err
	}
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
	ctx := context.Background()
	if err := sm.Redis.AddPlay(ctx, session.Game.ID, playResult); err != nil {
		return nil, err
	}
	if playResult.Outcome == game.OutcomeLose {
		if err := sm.finalizeGame(session); err != nil {
			return nil, err
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
	return sm.finalizeGame(session)
}

func (sm *SessionManager) finalizeGame(session *Session) error {
	ctx := context.Background()
	session.Balance += session.Game.CurrentBet
	if err := sm.saveGameToDB(session); err != nil {
		return err
	}
	if err := sm.Redis.DeleteGame(ctx, session.Game.ClientID); err != nil {
		return err
	}
	if err := sm.Redis.DeletePlays(ctx, session.Game.ID); err != nil {
		return err
	}
	if err := sm.DB.UpdateUserBalance(ctx, session.ClientID, session.Balance); err != nil {
		return err
	}
	session.Game = nil
	return nil
}

func (sm *SessionManager) saveGameToDB(session *Session) error {
	ctx := context.Background()
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
	dbgame := &db.Game{
		ID:         session.Game.ID,
		FinalBet:   session.Game.CurrentBet,
		TotalPlays: len(session.Game.Plays),
	}
	if err := sm.DB.UpdateGame(ctx, dbgame); err != nil {
		return err
	}
	return nil
}
