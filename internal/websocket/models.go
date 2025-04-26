package websocket

import "encoding/json"

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WalletRequest struct {
	ClientID string `json:"client_id"`
}

type WalletResponse struct {
	Balance int `json:"balance"`
}

type PlayRequest struct {
	ClientID  string `json:"client_id"`
	BetAmount int    `json:"bet_amount"`
	BetType   string `json:"bet_type"` // even / odd
}

type PlayResponse struct {
	RolledNumber int    `json:"rolled_number"`
	NextBet      int    `json:"next_bet"`
	Outcome      string `json:"outcome"` // win / lose
}

type EndPlayRequest struct {
	ClientID string `json:"client_id"`
}

type EndPlayResponse struct {
	Winnings int `json:"winnings"`
	Balance  int `json:"balance"`
}

type ReqWithClient interface {
	GetClientID() string
}

func (r WalletRequest) GetClientID() string {
	return r.ClientID
}

func (r PlayRequest) GetClientID() string {
	return r.ClientID
}

func (r EndPlayRequest) GetClientID() string {
	return r.ClientID
}
