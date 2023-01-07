package uno

import "github.com/Pato05/unobot/messages"

type GameAlreadyStartedError struct{}

func (self GameAlreadyStartedError) Error() string {
	return messages.GAME_ALREADY_STARTED_ERROR
}

type GameNotStartedError struct{}

func (self GameNotStartedError) Error() string {
	return messages.GAME_NOT_STARTED_ERROR
}

type CardNotPlayableError struct{}

func (self CardNotPlayableError) Error() string {
	return messages.CARD_NOT_PLAYABLE_ERROR
}

type PlayerAlreadyExistsError struct{}

func (self PlayerAlreadyExistsError) Error() string {
	return messages.PLAYER_ALREADY_IN_GAME_ERROR
}

type PlayerNotInGameError struct{}

func (self PlayerNotInGameError) Error() string {
	return messages.PLAYER_NOT_IN_GAME_ERROR
}

type TooManyPlayersError struct{}

func (self TooManyPlayersError) Error() string {
	return messages.TOO_MANY_PLAYERS_ERROR
}

type TooFewPlayersError struct{}

func (self TooFewPlayersError) Error() string {
	return messages.TOO_FEW_PLAYERS_ERROR
}

type GameDisbandedNoPlayers struct{}

func (self GameDisbandedNoPlayers) Error() string {
	return messages.GAME_DISBANDED_NO_PLAYERS
}

type GameDisbandedLastPlayerWon struct{}

func (self GameDisbandedLastPlayerWon) Error() string {
	return messages.GAME_DISBANDED_NO_PLAYERS
}

type LobbyClosedError struct{}

func (self LobbyClosedError) Error() string {
	return messages.CANT_JOIN_LOBBY_CLOSED
}
