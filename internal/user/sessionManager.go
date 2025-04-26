package user

import (
	"sync"

	"github.com/Anacardo89/doubleOrNothingDice/internal/game"
)

type Session struct {
	ClientID string
	Balance  int
	Game     *game.Game
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) Create(clientID string) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session := &Session{
		ClientID: clientID,
		Balance:  100,
	}
	sm.sessions[clientID] = session
	return session
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
