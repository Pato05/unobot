package bot

import (
	"github.com/Pato05/unobot/uno"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PlayerGame struct {
	Game       *uno.Game[*UnoPlayer]
	UnoPlayer  *UnoPlayer
	GameChatId int64
}

type GameManager struct {
	games   map[int64]*uno.Game[*UnoPlayer]
	players map[int64]PlayerGame
}

func (self *GameManager) GetGame(chatId int64) (*uno.Game[*UnoPlayer], error) {
	val, ok := self.games[chatId]
	if !ok {
		return nil, NoGameError{}
	}
	return val, nil
}

func (self *GameManager) GetPlayerGame(userId int64) (PlayerGame, bool) {
	player, found := self.players[userId]
	return player, found
}
func (self *GameManager) assertGameDoesntExist(chatId int64) error {
	_, ok := self.games[chatId]
	if ok {
		return OngoingGameError{}
	}
	return nil
}

func (self *GameManager) NewGame(chatId int64, userId int64) error {
	if err := self.assertGameDoesntExist(chatId); err != nil {
		return err
	}
	self.games[chatId] = &uno.Game[*UnoPlayer]{
		GameCreatorUID: userId,
	}
	return nil
}

func (self *GameManager) DeleteGame(chatId int64) error {
	game, err := self.GetGame(chatId)
	if err != nil {
		return err
	}

	for _, player := range game.Players {
		delete(self.players, player.GetUID())
	}

	delete(self.games, chatId)
	return nil
}

func (self *GameManager) PlayerJoin(chatId int64, user *tgbotapi.User) error {
	game, err := self.GetGame(chatId)
	if err != nil {
		return err
	}

	if _, ok := self.GetPlayerGame(user.ID); ok {
		return PlayerAlreadyInOtherGameError{}
	}

	player := &UnoPlayer{
		uno.Player{
			Name: user.FirstName,
			Id:   user.ID,
		},
	}

	if err := game.JoinPlayer(player); err != nil {
		return err
	}

	self.players[user.ID] = PlayerGame{
		Game:       game,
		UnoPlayer:  player,
		GameChatId: chatId,
	}

	return nil
}

// removes a player from a game.
func (self *GameManager) PlayerLeave(chatId int64, userId int64) (*UnoPlayer, error) {
	game, err := self.GetGame(chatId)
	if err != nil {
		return nil, PlayerNotInGameError{}
	}

	playerGame, ok := self.GetPlayerGame(userId)
	if !ok {
		return nil, PlayerNotInGameError{}
	}

	err = game.LeavePlayer(playerGame.UnoPlayer)

	switch err.(type) {
	case uno.GameDisbandedLastPlayerWon:
	default:
		delete(self.players, userId)
		return nil, err
	}

	self.DeleteGame(chatId)
	return game.CurrentPlayer(), err
}

func (self *GameManager) GetPlayersInGame(chatId int64) ([]*UnoPlayer, error) {
	game, err := self.GetGame(chatId)
	if err != nil {
		return nil, err
	}
	return game.Players, nil
}
