package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Manager struct {
	db *sqlx.DB
}

func NewManager(db *sqlx.DB) *Manager {
	return &Manager{db: db}
}

// USERS

func (m *Manager) CreateUser(ctx context.Context, user *User) error {
	stmt, err := m.db.PrepareNamedContext(ctx, CreateUserQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.GetContext(ctx, &user.ID, user)
}

func (m *Manager) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := m.db.GetContext(ctx, &user, GetUserByIDQuery, id)
	return &user, err
}

func (m *Manager) GetUserByName(ctx context.Context, username string) (*User, error) {
	var user User
	err := m.db.GetContext(ctx, &user, GetUserByNameQuery, username)
	return &user, err
}

func (m *Manager) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := m.db.GetContext(ctx, &user, GetUserByEmailQuery, email)
	return &user, err
}

func (m *Manager) ActivateUser(ctx context.Context, id string) error {
	_, err := m.db.ExecContext(ctx, ActivateUserQuery, id)
	return err
}

func (m *Manager) UpdateUserPassword(ctx context.Context, userID string, newHash string) error {
	_, err := m.db.ExecContext(ctx, UpdateUserPasswordQuery, newHash, userID)
	return err
}

func (m *Manager) IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := m.db.GetContext(ctx, &exists, CheckUsernameExistsQuery, username)
	return exists, err
}

func (m *Manager) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := m.db.GetContext(ctx, &exists, CheckEmailExistsQuery, email)
	return exists, err
}

// GAMES

func (m *Manager) CreateGame(ctx context.Context, game *Game) error {
	stmt, err := m.db.PrepareNamedContext(ctx, CreateGameQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.GetContext(ctx, &game.ID, game)
}

func (m *Manager) GetGameByID(ctx context.Context, id string) (*Game, error) {
	var game Game
	err := m.db.GetContext(ctx, &game, GetGameByIDQuery, id)
	return &game, err
}

func (m *Manager) GetGamesByUserID(ctx context.Context, userID string) ([]Game, error) {
	var games []Game
	err := m.db.SelectContext(ctx, &games, GetGamesByUserQuery, userID)
	return games, err
}

// PLAYS

func (m *Manager) CreatePlay(ctx context.Context, play *Play) error {
	stmt, err := m.db.PrepareNamedContext(ctx, CreatePlayQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.GetContext(ctx, &play.ID, play)
}

func (m *Manager) GetPlaysByGameID(ctx context.Context, gameID string) ([]Play, error) {
	var plays []Play
	err := m.db.SelectContext(ctx, &plays, GetPlaysByGameIDQuery, gameID)
	return plays, err
}
