# gounobot

A Telegram UNO bot written in Go

## Deployment

You need Go to compile the code base, so install it before proceeding.

First of all, please change the hardcoded userId inside constants/config.go

## Changing sticker pack

The bot uses the same sticker pack used by `@unobot` by default, but it can be changed by using the files inside the `gen/` folder.

First of all, uncomment `gen.HandleCardsGen()` in `bot.go` (but also comment `handler.ProcessUpdate()`), then proceed sending any message to the bot, and send the required stickers one by one.

After doing so, stop the bot by killing the process (`CTRL+C`), comment the line again, and run `gen/card_gen.py` to generate the other needed maps.

## Credits

-   [jh0ker/mau_mau_bot](https://github.com/jh0ker/mau_mau_bot), for the idea, as well as some inspiration for some parts of the code
