package messages

const (
	CARD_NOT_PLAYABLE_ERROR              = "Card not playable"
	PLAYER_ALREADY_IN_GAME_ERROR         = "Player is already in the game!"
	ONGOING_GAME_ERROR                   = "There's already an ongoing game!"
	NO_GAMES_ERROR                       = "No games found."
	GAME_CREATED_SUCCESS                 = "Game created! Use /gojoin to join, then /gostart to start the game when you're ready!"
	GAME_JOINED_SUCCESS                  = "Joined game successfully"
	GAME_LEAVE_SUCCESS                   = "You left the game."
	NO_PLAYERS_IN_GAME_ERROR             = "No players in the game!"
	GAME_ALREADY_STARTED_ERROR           = "Game already started!"
	GAME_NOT_STARTED_ERROR               = "The game hasn't started yet!"
	PLAYER_ALREADY_IN_ANOTHER_GAME_ERROR = "You've joined another game already!"
	PLAYER_NOT_IN_GAME_ERROR             = "You're not playing in a game!"
	TOO_MANY_PLAYERS_ERROR               = "This game has reached the maximum amount of players!"
	TOO_FEW_PLAYERS_ERROR                = "Can't start game, at least two players must join the game."
	LOBBY_CLOSED_SUCCESS                 = "The lobby has been closed."
	LOBBY_OPEN_SUCCESS                   = "The lobby is now open to new players!"
	GAME_DISBANDED_NO_PLAYERS            = "The game was disbanded, not enough players."
	GAME_DISBANDED_LAST_PLAYER_WON       = "The game was disbanded, the last player won."
	CANT_JOIN_LOBBY_CLOSED               = "Can't join the game. The lobby is closed."
)
