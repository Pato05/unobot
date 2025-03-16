package bot

import "github.com/Pato05/unobot/messages"

type OngoingGameError struct{}

func (e OngoingGameError) Error() string {
	return messages.ONGOING_GAME_ERROR
}

type NoGameError struct{}

func (e NoGameError) Error() string {
	return messages.NO_GAMES_ERROR
}

type PlayerAlreadyInOtherGameError struct{}

func (e PlayerAlreadyInOtherGameError) Error() string {
	return messages.PLAYER_ALREADY_IN_ANOTHER_GAME_ERROR
}

type PlayerNotInGameError struct{}

func (e PlayerNotInGameError) Error() string {
	return messages.PLAYER_NOT_IN_GAME_ERROR
}
