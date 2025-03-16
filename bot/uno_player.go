package bot

import (
	"html"
	"strconv"

	"github.com/Pato05/unobot/uno"
)

type UnoPlayer struct {
	uno.Player
}

func (p *UnoPlayer) HTML() string {
	return "<a href=\"tg://user?id=" + strconv.Itoa(int(p.Id)) + "\">" + html.EscapeString(p.Name) + "</a>"
}

func (p *UnoPlayer) EscapedName() string {
	return html.EscapeString(p.Name)
}
