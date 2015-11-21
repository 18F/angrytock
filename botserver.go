/*
Script for starting the bot server
*/

package main

import (
	"net/http"
	"os"

	"github.com/geramirez/tock-bot/bot"
	"github.com/robfig/cron"
)

func main() {

	bot := bot.InitBot()

	// Slap tock users a couple times
	c := cron.New()
	c.AddFunc("@weekly", func() {
		bot.StoreSlackUsers()
	})
	c.AddFunc("0 0 0 * * 1", func() {
		bot.SlapLateUsers()
	})
	c.Start()

	bot.BotherSlackUsers()

	// Start server
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}
