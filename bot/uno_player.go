package bot

import (
	"html"
	"strconv"

	"github.com/Pato05/unobot/uno"
)

type UnoPlayer struct {
	uno.Player
}

func (self *UnoPlayer) HTML() string {
	return "<a href=\"tg://user?id=" + strconv.Itoa(int(self.Id)) + "\">" + html.EscapeString(self.Name) + "</a>"
}

func (self *UnoPlayer) EscapedName() string {
	return html.EscapeString(self.Name)
}
