package db

const (
	// Users
	CreateUserQuery = `
        INSERT INTO users (
			username, 
			email, 
			password_hash
		)
        VALUES (
			:username, 
			:email, 
			:password_hash
		)
		RETURNING id
    ;`

	GetUserByIDQuery = `
        SELECT *
        FROM users
        WHERE id = $1
    ;`

	GetUserByEmailQuery = `
        SELECT *
        FROM users
        WHERE email = $1
    ;`

	ActivateUserQuery = `
		UPDATE users
		SET is_active = TRUE
		WHERE id = $1
	`

	// Games
	CreateGameQuery = `
        INSERT INTO games (
			user_id, 
			initial_bet, 
			final_bet, 
			total_plays
		)
        VALUES (
			:user_id, 
			:initial_bet, 
			:final_bet, 
			:total_plays
		)
		RETURNING id
    ;`

	GetGameByIDQuery = `
        SELECT *
        FROM games
        WHERE id = $1
    ;`

	GetGamesByUserQuery = `
        SELECT *
        FROM games
        WHERE user_id = $1
    ;`

	// Plays
	CreatePlayQuery = `
        INSERT INTO plays (
			game_id, 
			play_number,
			bet_amount, 
			play_choice, 
			dice_result, 
			outcome
		)
        VALUES (
			:game_id, 
			:play_number, 
			:bet_amount, 
			:play_choice, 
			:dice_result, 
			:outcome
		)
		RETURNING id
    ;`

	GetPlaysByGameIDQuery = `
        SELECT *
        FROM plays
        WHERE game_id = $1
        ORDER BY play_number ASC
    ;`
)
