package db

import "time"

type User struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Game struct {
	ID         string     `db:"id"`
	UserID     string     `db:"user_id"`
	InitialBet int        `db:"initial_bet"`
	FinalBet   int        `db:"final_bet"`
	TotalPlays int        `db:"total_plays"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	EndTime    *time.Time `db:"end_time"`
}

type Play struct {
	ID         string    `db:"id"`
	GameID     string    `db:"game_id"`
	PlayNumber int       `db:"play_number"`
	BetAmount  int       `db:"bet_amount"`
	PlayChoice string    `db:"play_choice"`
	DiceResult int       `db:"dice_result"`
	Outcome    string    `db:"outcome"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
