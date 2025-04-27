CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    balance INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    initial_bet INTEGER NOT NULL,
    final_bet INTEGER NOT NULL,
    total_plays INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    end_time TIMESTAMP DEFAULT NULL
);

CREATE TABLE plays (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID REFERENCES games(id),
    play_number INTEGER NOT NULL,
    bet_amount INTEGER NOT NULL,
    play_choice VARCHAR(10) NOT NULL,
    dice_result INTEGER NOT NULL,
    outcome VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_games_updated_at
BEFORE UPDATE ON games
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_plays_updated_at
BEFORE UPDATE ON plays
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();