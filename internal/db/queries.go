package db

const (
	// Users
	CreateUserQuery = `
        INSERT INTO users (
			username, 
			email, 
			password_hash,
			balance
		)
        VALUES (
			:username, 
			:email, 
			:password_hash,
			0
		)
		RETURNING id
    ;`

	GetUserByIDQuery = `
        SELECT *
        FROM users
        WHERE id = $1
    ;`

	GetUserByNameQuery = `
        SELECT *
        FROM users
        WHERE username = $1
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
	;`

	UpdateUserPasswordQuery = `
		UPDATE users 
		SET password_hash = $1 
		WHERE id = $2
	;`

	UpdateUserBalanceQuery = `
		UPDATE users 
		SET balance = $1 
		WHERE id = $2
	;`

	CheckUsernameExistsQuery = `
		SELECT EXISTS (
			SELECT 1 FROM users 
			WHERE username = $1
		)
	;`

	CheckEmailExistsQuery = `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE email = $1
		)
	;`

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

	UpdateGameQuery = `
		UPDATE games
		SET
    		final_bet = :final_bet,
    		total_plays = :total_plays,
    		end_time = NOW()
		WHERE
    		id = :game_id
		RETURNING id;
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
