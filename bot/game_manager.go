package bot

import (
	"time"

	"github.com/Pato05/unobot/uno"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PlayerGame struct {
	Game      *UnoGame
	UnoPlayer *UnoPlayer
	LastTimer *time.Timer
}

type GameManager struct {
	games   map[int64]*UnoGame
	players map[int64]PlayerGame
}

func (gm *GameManager) GetGame(chatId int64) (*UnoGame, error) {
	val, ok := gm.games[chatId]
	if !ok {
		return nil, NoGameError{}
	}
	return val, nil
}

func (gm *GameManager) GetPlayerGame(userId int64) (*PlayerGame, bool) {
	player, found := gm.players[userId]
	return &player, found
}

func (gm *GameManager) assertGameDoesntExist(chatId int64) error {
	_, ok := gm.games[chatId]
	if ok {
		return OngoingGameError{}
	}
	return nil
}

func (gm *GameManager) NewGame(chatId int64, userId int64) error {
	if err := gm.assertGameDoesntExist(chatId); err != nil {
		return err
	}

	gm.games[chatId] = &UnoGame{
		Game: uno.Game[*UnoPlayer]{
			GameCreatorUID: userId,
		},
		ChatId: chatId,
	}

	return nil
}

func (gm *GameManager) DeleteGame(chatId int64) error {
	game, err := gm.GetGame(chatId)
	if err != nil {
		return err
	}

	for _, player := range game.Players {
		delete(gm.players, player.GetUID())
	}

	delete(gm.games, chatId)
	return nil
}

func (gm *GameManager) PlayerJoin(chatId int64, user *tgbotapi.User) error {
	if _, ok := gm.GetPlayerGame(user.ID); ok {
		return PlayerAlreadyInOtherGameError{}
	}

	game, err := gm.GetGame(chatId)
	if err != nil {
		return err
	}

	player := &UnoPlayer{
		Player: uno.Player{
			Id: user.ID,
		},

		Name: user.FirstName,
	}

	if err := game.JoinPlayer(player); err != nil {
		return err
	}

	gm.players[user.ID] = PlayerGame{
		Game:      game,
		UnoPlayer: player,
	}

	return nil
}

// removes a player from a game.
func (gm *GameManager) PlayerLeave(chatId int64, userId int64) (*UnoPlayer, error) {
	game, err := gm.GetGame(chatId)
	if err != nil {
		return nil, PlayerNotInGameError{}
	}

	playerGame, ok := gm.GetPlayerGame(userId)
	if !ok {
		return nil, PlayerNotInGameError{}
	}

	err = game.LeavePlayer(playerGame.UnoPlayer)

	switch err.(type) {
	case uno.GameDisbandedLastPlayerWon:
	default:
		delete(gm.players, userId)
		return nil, err
	}

	gm.DeleteGame(chatId)
	return game.CurrentPlayer(), err
}

func (gm *GameManager) GetPlayersInGame(chatId int64) ([]*UnoPlayer, error) {
	game, err := gm.GetGame(chatId)
	if err != nil {
		return nil, err
	}
	return game.Players, nil
}
