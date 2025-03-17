package bot

import (
	"html"
	"strconv"
	"time"

	"github.com/Pato05/unobot/constants"
	"github.com/Pato05/unobot/uno"
)

type UnoPlayer struct {
	uno.Player

	Name string
}

func (p *UnoPlayer) HTML() string {
	return "<a href=\"tg://user?id=" + strconv.Itoa(int(p.Id)) + "\">" + html.EscapeString(p.Name) + "</a>"
}

func (p *UnoPlayer) EscapedName() string {
	return html.EscapeString(p.Name)
}

type UnoGame struct {
	uno.Game[*UnoPlayer]

	LastTimer *time.Timer
	ChatId    int64
}

func (g *UnoGame) StopCurrentTimer() {
	if g.LastTimer != nil {
		g.LastTimer.Stop()
	}
}

func (g *UnoGame) SetTimer(bh *BotHandler) (bool, error) {
	g.StopCurrentTimer()

	// set the timer for the current player
	player := g.CurrentPlayer()
	if player.AutoSkipCount >= constants.KICK_AFTER_SKIP {
		// kick current player
		bh.announceKickedAFKPlayer(g.ChatId, player)
		delete(bh.gameManager.players, player.GetUID())
		if err := g.LeaveCurrentPlayer(); err != nil {
			bh.SendMessage(g.ChatId, err.Error())
			bh.gameManager.DeleteGame(g.ChatId)
			return false, err
		}

		return g.SetTimer(bh)
	}

	timeout := player.SkipTimer()
	g.LastTimer = time.AfterFunc(timeout, func() {
		bh.logDebug("Skipping player...")
		// increase the auto skip count and skip this player
		player.IncreaseAutoSkipCount()
		g.PreAutoSkipPlayer()
		bh.announcePlayerSkipped(g.ChatId, player)
		bh.nextPlayer(g)
	})

	bh.logDebug("timer is set")

	return true, nil
}
