package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Anacardo89/doubleOrNothingDice/internal/game"

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	client *redis.Client
}

func NewManager(client *redis.Client) *Manager {
	return &Manager{client: client}
}

// Game

func (rm *Manager) SaveGame(ctx context.Context, g *game.Game) error {
	key := g.ClientID
	gameMap := map[string]interface{}{
		"ID":         g.ID,
		"InitialBet": g.InitialBet,
		"CurrentBet": g.CurrentBet,
		"IsActive":   g.IsActive,
	}
	return rm.client.HSet(ctx, key, gameMap).Err()
}

func (rm *Manager) GetGame(ctx context.Context, userID string) (*game.Game, error) {
	key := userID
	data, err := rm.client.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("no game found for user %s", userID)
	}
	initialBet, _ := strconv.Atoi(data["InitialBet"])
	currentBet, _ := strconv.Atoi(data["CurrentBet"])
	isActive, _ := strconv.ParseBool(data["IsActive"])
	return &game.Game{
		ID:         data["ID"],
		ClientID:   key,
		InitialBet: initialBet,
		CurrentBet: currentBet,
		IsActive:   isActive,
	}, nil
}

func (rm *Manager) DeleteGame(ctx context.Context, userID string) error {
	return rm.client.Del(ctx, userID).Err()
}

// Play

func (rm *Manager) AddPlay(ctx context.Context, gameID string, play *game.Play) error {
	data, err := json.Marshal(play)
	if err != nil {
		return fmt.Errorf("failed to marshal play: %w", err)
	}
	return rm.client.RPush(ctx, gameID, data).Err()
}

func (rm *Manager) GetPlays(ctx context.Context, gameID string) ([]game.Play, error) {
	raw, err := rm.client.LRange(ctx, gameID, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	plays := make([]game.Play, len(raw))
	for i, item := range raw {
		if err := json.Unmarshal([]byte(item), &plays[i]); err != nil {
			return nil, fmt.Errorf("error decoding play %d: %w", i, err)
		}
	}
	return plays, nil
}

func (rm *Manager) DeletePlays(ctx context.Context, gameID string) error {
	return rm.client.Del(ctx, gameID).Err()
}

// Abstract

func (m *Manager) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(ctx, key, jsonData, ttl).Err()
}

func (m *Manager) GetJSON(ctx context.Context, key string, dest any) error {
	data, err := m.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
